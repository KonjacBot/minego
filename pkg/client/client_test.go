package client

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"
	protocolpacket "github.com/KonjacBot/minego/pkg/protocol/packet"
	configserver "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/server"
	gameclient "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

func TestNewClientAllowsNilOptions(t *testing.T) {
	c := NewClient(nil)
	if c == nil {
		t.Fatal("NewClient(nil) returned nil")
	}
}

func TestCloseBeforeConnect(t *testing.T) {
	c := NewClient(nil)
	if err := c.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestConfigurationReturnsReadError(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	b := &botClient{conn: mcnet.WrapConn(clientConn)}
	_ = serverConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := b.configuration(ctx); err == nil {
		t.Fatal("configuration() returned nil after peer closed")
	}
}

func TestHandleGameStopsPromptlyWhenContextIsCanceled(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()
	b := &botClient{
		conn:          mcnet.WrapConn(clientConn),
		packetHandler: newPacketHandler(),
		eventHandler:  NewEventHandler(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- b.HandleGame(ctx) }()
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("HandleGame() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("HandleGame did not stop after context cancellation")
	}
}

func TestWritePacketRechecksContextAfterWaitingForLock(t *testing.T) {
	conn := &countingConn{}
	b := &botClient{conn: mcnet.WrapConn(conn)}
	b.writeMu.Lock()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- b.WritePacket(ctx, &server.Pong{ID: 1})
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	b.writeMu.Unlock()

	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("WritePacket() error = %v, want context.Canceled", err)
	}
	if writes := conn.writes.Load(); writes != 0 {
		t.Fatalf("canceled WritePacket performed %d writes", writes)
	}
}

func TestWritePacketReturnsEncodingErrors(t *testing.T) {
	conn := &countingConn{}
	b := &botClient{conn: mcnet.WrapConn(conn)}

	err := b.WritePacket(context.Background(), &server.Chat{Message: "hello"})
	if err == nil {
		t.Fatal("WritePacket() accepted a chat packet without its fixed bit set")
	}
	if writes := conn.writes.Load(); writes != 0 {
		t.Fatalf("invalid packet performed %d writes", writes)
	}
}

func TestWritePacketCancellationInterruptsInFlightWrite(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()
	b := &botClient{conn: mcnet.WrapConn(clientConn)}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- b.WritePacket(ctx, &server.Pong{ID: 1}) }()
	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("WritePacket() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("WritePacket() did not stop after context cancellation")
	}
}

func TestCloseInterruptsBlockedWrite(t *testing.T) {
	conn := newBlockingConn()
	b := &botClient{conn: mcnet.WrapConn(conn)}
	writeDone := make(chan error, 1)
	go func() { writeDone <- b.WritePacket(context.Background(), &server.Pong{ID: 1}) }()

	select {
	case <-conn.started:
	case <-time.After(time.Second):
		t.Fatal("write did not start")
	}
	if err := b.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	select {
	case err := <-writeDone:
		if err == nil {
			t.Fatal("blocked write returned nil after close")
		}
	case <-time.After(time.Second):
		t.Fatal("Close() did not interrupt the blocked write")
	}
}

func TestHandleGameContainsRawHandlerPanic(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()
	b := &botClient{
		conn:          mcnet.WrapConn(clientConn),
		packetHandler: newPacketHandler(),
		eventHandler:  NewEventHandler(),
	}
	b.packetHandler.AddRawPacketHandler(packetid.ClientboundKeepAlive, func(context.Context, pk.Packet) {
		panic("raw handler failure")
	})
	done := make(chan error, 1)
	go func() { done <- b.HandleGame(context.Background()) }()
	peer := mcnet.WrapConn(serverConn)
	if err := peer.WritePacket(pk.Marshal(packetid.ClientboundKeepAlive, pk.Long(1))); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-done:
		var panicErr *CallbackPanicError
		if !errors.As(err, &panicErr) {
			t.Fatalf("HandleGame() error = %v, want CallbackPanicError", err)
		}
	case <-time.After(time.Second):
		t.Fatal("HandleGame() did not report raw handler panic")
	}
}

func TestHandleGameContainsTypedHandlerPanic(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()
	b := &botClient{
		conn:          mcnet.WrapConn(clientConn),
		packetHandler: newPacketHandler(),
		eventHandler:  NewEventHandler(),
	}
	b.packetHandler.AddPacketHandler(packetid.ClientboundKeepAlive, func(context.Context, gameclient.ClientboundPacket) {
		panic("typed handler failure")
	})
	done := make(chan error, 1)
	go func() { done <- b.HandleGame(context.Background()) }()
	peer := mcnet.WrapConn(serverConn)
	if err := peer.WritePacket(pk.Marshal(packetid.ClientboundKeepAlive, pk.Long(1))); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-done:
		var panicErr *CallbackPanicError
		if !errors.As(err, &panicErr) {
			t.Fatalf("HandleGame() error = %v, want CallbackPanicError", err)
		}
	case <-time.After(time.Second):
		t.Fatal("HandleGame() did not report typed handler panic")
	}
}

func TestConfigurationHandlesResourcePackAndCodeOfConduct(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	acceptedConn := make(chan net.Conn, 1)
	acceptErr := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			acceptErr <- err
			return
		}
		acceptedConn <- conn
	}()

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer clientConn.Close()

	var serverConn net.Conn
	select {
	case serverConn = <-acceptedConn:
		defer serverConn.Close()
	case err = <-acceptErr:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for test server accept")
	}

	b := &botClient{conn: mcnet.WrapConn(clientConn)}
	peer := mcnet.WrapConn(serverConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- b.configuration(ctx) }()
	var response pk.Packet
	if err := peer.ReadPacket(&response); err != nil {
		t.Fatal(err)
	}
	if got := packetid.ServerboundPacketID(response.ID); got != packetid.ServerboundConfigClientInformation {
		t.Fatalf("initial packet ID = %v, want %v", got, packetid.ServerboundConfigClientInformation)
	}
	var clientInformation configserver.ConfigClientInformation
	if err := protocolpacket.Scan(response, &clientInformation); err != nil {
		t.Fatalf("decode client information: %v", err)
	}

	packID := uuid.MustParse("12345678-1234-5678-9abc-def012345678")
	if err := peer.WritePacket(pk.Marshal(
		packetid.ClientboundConfigResourcePackPush,
		pk.UUID(packID),
		pk.String("https://example.invalid/pack.zip"),
		pk.String(strings.Repeat("0", 40)),
		pk.Boolean(false),
		pk.Boolean(false),
	)); err != nil {
		t.Fatal(err)
	}
	if err := peer.ReadPacket(&response); err != nil {
		t.Fatal(err)
	}
	if got := packetid.ServerboundPacketID(response.ID); got != packetid.ServerboundConfigResourcePack {
		t.Fatalf("response packet ID = %v, want %v", got, packetid.ServerboundConfigResourcePack)
	}
	var gotPackID pk.UUID
	var result pk.VarInt
	if err := response.Scan(&gotPackID, &result); err != nil {
		t.Fatal(err)
	}
	if uuid.UUID(gotPackID) != packID || int32(result) != resourcePackResultDeclined {
		t.Fatalf("resource pack response = (%v, %d)", uuid.UUID(gotPackID), result)
	}

	if err := peer.WritePacket(pk.Marshal(packetid.ClientboundConfigCodeOfConduct, pk.String("be nice"))); err != nil {
		t.Fatal(err)
	}
	if err := peer.ReadPacket(&response); err != nil {
		t.Fatal(err)
	}
	if got := packetid.ServerboundPacketID(response.ID); got != packetid.ServerboundConfigAcceptCodeOfConduct {
		t.Fatalf("response packet ID = %v, want %v", got, packetid.ServerboundConfigAcceptCodeOfConduct)
	}
	if len(response.Data) != 0 {
		t.Fatalf("accept code of conduct payload = %v, want empty", response.Data)
	}

	if err := peer.WritePacket(pk.Marshal(packetid.ClientboundConfigFinishConfiguration)); err != nil {
		t.Fatal(err)
	}
	if err := peer.ReadPacket(&response); err != nil {
		t.Fatal(err)
	}
	if got := packetid.ServerboundPacketID(response.ID); got != packetid.ServerboundConfigFinishConfiguration {
		t.Fatalf("response packet ID = %v, want %v", got, packetid.ServerboundConfigFinishConfiguration)
	}

	if err := <-done; err != nil {
		t.Fatalf("configuration() error = %v", err)
	}
}

func TestDecodeClientboundPacketRejectsTrailingData(t *testing.T) {
	data := pk.Marshal(packetid.ClientboundKeepAlive, pk.Long(99)).Data
	data = append(data, 0x01)

	pkt, handled, err := decodeClientboundPacket(packetid.ClientboundKeepAlive, data)
	if pkt != nil || !handled || err == nil {
		t.Fatalf("decodeClientboundPacket() = (%T, %t, %v), want handled trailing-data error", pkt, handled, err)
	}
}

func TestDecodeClientboundPacketUsesReaderConsumption(t *testing.T) {
	if err := decodeClientboundPayload(&underreportingPacket{}, []byte{0x2a}); err != nil {
		t.Fatalf("decodeClientboundPayload() error = %v, want successful consumed decode", err)
	}
}

func TestDecodeClientboundPacketSkipsUnsupportedPacket(t *testing.T) {
	pkt, handled, err := decodeClientboundPacket(packetid.ClientboundPlayerChat, []byte{0x01, 0x02})
	if pkt != nil || handled || err != nil {
		t.Fatalf("decodeClientboundPacket() = (%T, %t, %v), want unsupported skip", pkt, handled, err)
	}
}

func TestDecodeClientboundPacketUsesFixedResourcePackRegistry(t *testing.T) {
	packID := uuid.MustParse("12345678-1234-5678-9abc-def012345678")
	data := pk.Marshal(
		packetid.ClientboundResourcePackPush,
		pk.UUID(packID),
		pk.String("https://example.invalid/pack.zip"),
		pk.String(strings.Repeat("0", 40)),
		pk.Boolean(false),
		pk.Boolean(false),
	).Data

	pkt, handled, err := decodeClientboundPacket(packetid.ClientboundResourcePackPush, data)
	if err != nil {
		t.Fatal(err)
	}
	if !handled {
		t.Fatal("decodeClientboundPacket() skipped resource pack push")
	}
	resourcePack, ok := pkt.(*gameclient.AddResourcePack)
	if !ok {
		t.Fatalf("decoded packet = %T, want *gameclient.AddResourcePack", pkt)
	}
	if resourcePack.UUID != packID {
		t.Fatalf("decoded resource pack UUID = %v, want %v", resourcePack.UUID, packID)
	}
}

func TestDecodeClientboundPacketRejectsNegativeStringLength(t *testing.T) {
	data := make([]byte, 16)
	data = append(data, 0xff, 0xff, 0xff, 0xff, 0x0f)
	packet, handled, err := decodeClientboundPacket(packetid.ClientboundResourcePackPush, data)
	if packet != nil || !handled || err == nil {
		t.Fatalf("decodeClientboundPacket() = (%T, %t, %v), want handled decode error", packet, handled, err)
	}
}

type countingConn struct {
	writes atomic.Int32
}

type blockingConn struct {
	started chan struct{}
	closed  chan struct{}
	start   sync.Once
	close   sync.Once
}

func newBlockingConn() *blockingConn {
	return &blockingConn{started: make(chan struct{}), closed: make(chan struct{})}
}

type underreportingPacket struct{}

func (*underreportingPacket) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundKeepAlive
}

func (*underreportingPacket) ReadFrom(r io.Reader) (int64, error) {
	var data [1]byte
	_, err := io.ReadFull(r, data[:])
	return 0, err
}

func (*underreportingPacket) WriteTo(io.Writer) (int64, error) { return 0, nil }

func (*countingConn) Read([]byte) (int, error)         { return 0, io.EOF }
func (c *countingConn) Write(p []byte) (int, error)    { c.writes.Add(1); return len(p), nil }
func (*countingConn) Close() error                     { return nil }
func (*countingConn) LocalAddr() net.Addr              { return testAddr("local") }
func (*countingConn) RemoteAddr() net.Addr             { return testAddr("remote") }
func (*countingConn) SetDeadline(time.Time) error      { return nil }
func (*countingConn) SetReadDeadline(time.Time) error  { return nil }
func (*countingConn) SetWriteDeadline(time.Time) error { return nil }

func (c *blockingConn) Read([]byte) (int, error) { return 0, io.EOF }
func (c *blockingConn) Write([]byte) (int, error) {
	c.start.Do(func() { close(c.started) })
	<-c.closed
	return 0, io.ErrClosedPipe
}
func (c *blockingConn) Close() error {
	c.close.Do(func() { close(c.closed) })
	return nil
}
func (*blockingConn) LocalAddr() net.Addr              { return testAddr("local") }
func (*blockingConn) RemoteAddr() net.Addr             { return testAddr("remote") }
func (*blockingConn) SetDeadline(time.Time) error      { return nil }
func (*blockingConn) SetReadDeadline(time.Time) error  { return nil }
func (*blockingConn) SetWriteDeadline(time.Time) error { return nil }

type testAddr string

func (a testAddr) Network() string { return string(a) }
func (a testAddr) String() string  { return string(a) }
