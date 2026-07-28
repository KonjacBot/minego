package client

import (
	"context"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	mcnet "github.com/KonjacBot/go-mc/net"
	"github.com/KonjacBot/minego/pkg/auth"
	"github.com/KonjacBot/minego/pkg/bot"
	configclient "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

func TestNewClientAllowsNilOptions(t *testing.T) {
	c := NewClient(nil)
	if c == nil {
		t.Fatal("NewClient(nil) returned nil")
	}
}

func TestNewClientUsesConfiguredReadIdleTimeout(t *testing.T) {
	const timeout = 17 * time.Second
	c := NewClient(&bot.ClientOptions{ReadIdleTimeout: timeout}).(*botClient)
	if c.readIdleTimeout != timeout {
		t.Fatalf("read idle timeout = %s, want %s", c.readIdleTimeout, timeout)
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

func TestSupportedKnownPacksSelectsOnlyMatchingCorePack(t *testing.T) {
	offered := []configclient.KnownPack{
		{Namespace: "example", ID: "custom", Version: "1"},
		{Namespace: "minecraft", ID: "core", Version: "26.1"},
		{Namespace: "minecraft", ID: "core", Version: "26.2"},
	}
	selected := supportedKnownPacks(offered)
	if len(selected) != 1 || selected[0] != offered[2] {
		t.Fatalf("selected packs = %+v", selected)
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

func TestHandleGameReturnsAfterReadIdleTimeout(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()
	b := &botClient{
		conn:            mcnet.WrapConn(clientConn),
		packetHandler:   newPacketHandler(),
		eventHandler:    NewEventHandler(),
		readIdleTimeout: 25 * time.Millisecond,
	}

	started := time.Now()
	err := b.HandleGame(context.Background())
	if !isNetworkTimeout(err) {
		t.Fatalf("HandleGame() error = %v, want network timeout", err)
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("HandleGame() returned after %s, want prompt idle timeout", elapsed)
	}
}

func TestConfigurationReturnsAfterReadIdleTimeout(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()
	b := &botClient{
		conn:            mcnet.WrapConn(clientConn),
		readIdleTimeout: 25 * time.Millisecond,
	}

	err := b.configuration(context.Background())
	if !isNetworkTimeout(err) {
		t.Fatalf("configuration() error = %v, want network timeout", err)
	}
}

func TestLoginReturnsAfterReadIdleTimeout(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer serverConn.Close()
	go func() {
		_, _ = io.Copy(io.Discard, serverConn)
	}()
	b := &botClient{
		conn:            mcnet.WrapConn(clientConn),
		authProvider:    &auth.OfflineAuth{Username: "idle-test"},
		readIdleTimeout: 25 * time.Millisecond,
	}

	err := b.login(context.Background())
	if !isNetworkTimeout(err) {
		t.Fatalf("login() error = %v, want network timeout", err)
	}
}

func isNetworkTimeout(err error) bool {
	var networkError net.Error
	return errors.As(err, &networkError) && networkError.Timeout()
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

type countingConn struct {
	writes atomic.Int32
}

func (*countingConn) Read([]byte) (int, error)         { return 0, io.EOF }
func (c *countingConn) Write(p []byte) (int, error)    { c.writes.Add(1); return len(p), nil }
func (*countingConn) Close() error                     { return nil }
func (*countingConn) LocalAddr() net.Addr              { return testAddr("local") }
func (*countingConn) RemoteAddr() net.Addr             { return testAddr("remote") }
func (*countingConn) SetDeadline(time.Time) error      { return nil }
func (*countingConn) SetReadDeadline(time.Time) error  { return nil }
func (*countingConn) SetWriteDeadline(time.Time) error { return nil }

type testAddr string

func (a testAddr) Network() string { return string(a) }
func (a testAddr) String() string  { return string(a) }
