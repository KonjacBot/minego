package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/v2/pkg/protocol"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet"

	"github.com/KonjacBot/minego/v2/pkg/auth"
	"github.com/KonjacBot/minego/v2/pkg/bot"
	"github.com/KonjacBot/minego/v2/pkg/bus"
	"github.com/KonjacBot/minego/v2/pkg/game/inventory"
	"github.com/KonjacBot/minego/v2/pkg/game/player"
	"github.com/KonjacBot/minego/v2/pkg/game/world"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
)

type Client struct {
	conn               *mcnet.Conn
	connMu             sync.RWMutex
	writeMu            sync.Mutex
	stateMu            sync.Mutex
	state              packet.State
	stateChanged       chan struct{}
	bus                *bus.Bus
	connected          atomic.Bool
	authProvider       auth.Provider
	readIdleTimeout    time.Duration
	resourcePackPolicy bot.ResourcePackPolicy
	cookieMu           sync.RWMutex
	cookies            map[string][]int8
	registryMu         sync.RWMutex
	registryEntries    map[string]map[int32]string

	inventory *inventory.Manager
	world     *world.World
	player    *player.Player
}

func (b *Client) Player() bot.Player {
	return b.player
}

func (b *Client) Close(ctx context.Context) error {
	b.writeMu.Lock()
	defer b.writeMu.Unlock()

	b.connMu.Lock()
	conn := b.conn
	b.conn = nil
	b.connected.Store(false)
	b.connMu.Unlock()
	b.registryMu.Lock()
	b.registryEntries = nil
	b.registryMu.Unlock()
	b.setConnectionStateLocked(packet.StateLogin)

	if conn == nil {
		return ctx.Err()
	}
	if err := conn.Close(); err != nil {
		return err
	}
	return ctx.Err()
}

func (b *Client) IsConnected() bool {
	return b.connected.Load()
}

func (b *Client) WritePacket(ctx context.Context, outbound server.ServerboundPacket) error {
	raw := pk.Marshal(outbound.PacketID(), outbound)
	for {
		if err := b.waitForPlayState(ctx); err != nil {
			return err
		}

		b.writeMu.Lock()
		if err := ctx.Err(); err != nil {
			b.writeMu.Unlock()
			return err
		}
		if b.connectionState() != packet.StatePlay {
			b.writeMu.Unlock()
			continue
		}
		err := b.writeRawPacketLocked(ctx, raw)
		b.writeMu.Unlock()
		return err
	}
}

func (b *Client) writeRawPacket(ctx context.Context, packet pk.Packet) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.writeMu.Lock()
	defer b.writeMu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}

	return b.writeRawPacketLocked(ctx, packet)
}

func (b *Client) writeRawPacketLocked(ctx context.Context, packet pk.Packet) error {
	b.connMu.RLock()
	conn := b.conn
	b.connMu.RUnlock()
	if conn == nil {
		return errors.New("client is not connected")
	}

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.Socket.SetWriteDeadline(deadline); err != nil {
			return err
		}
		defer conn.Socket.SetWriteDeadline(time.Time{})
	}

	return conn.WritePacket(packet)
}

func (b *Client) waitForPlayState(ctx context.Context) error {
	for {
		b.connMu.RLock()
		conn := b.conn
		b.connMu.RUnlock()
		if conn == nil {
			return errors.New("client is not connected")
		}

		b.stateMu.Lock()
		if b.state == packet.StatePlay {
			b.stateMu.Unlock()
			return nil
		}
		if b.stateChanged == nil {
			b.stateChanged = make(chan struct{})
		}
		changed := b.stateChanged
		b.stateMu.Unlock()

		b.connMu.RLock()
		conn = b.conn
		b.connMu.RUnlock()
		if conn == nil {
			return errors.New("client is not connected")
		}

		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case <-changed:
		}
	}
}

func (b *Client) connectionState() packet.State {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	return b.state
}

func (b *Client) setConnectionState(state packet.State) {
	b.writeMu.Lock()
	defer b.writeMu.Unlock()
	b.setConnectionStateLocked(state)
}

func (b *Client) setConnectionStateLocked(state packet.State) {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	b.state = state
	if b.stateChanged != nil {
		close(b.stateChanged)
	}
	b.stateChanged = make(chan struct{})
}

// Bus returns the bus that delivers this client's packets and events. The
// methods below forward to it, so most code can call them on the Client.
func (b *Client) Bus() *bus.Bus {
	return b.bus
}

// On subscribes f to a packet type (e.g. *client.KeepAlive) or an event type
// (e.g. player.MessageEvent) and returns a function that removes it.
func (b *Client) On[T any](f func(ctx context.Context, v T)) (off func()) {
	return b.bus.On[T](f)
}

// OnCtx is On, removed automatically when ctx is done.
func (b *Client) OnCtx[T any](ctx context.Context, f func(ctx context.Context, v T)) (off func()) {
	return b.bus.OnCtx[T](ctx, f)
}

// Once is On for a single delivery.
func (b *Client) Once[T any](f func(ctx context.Context, v T)) (off func()) {
	return b.bus.Once[T](f)
}

// WaitFor blocks until a packet or event of type T matches, or ctx is done.
func (b *Client) WaitFor[T any](ctx context.Context, match func(v T) bool) (T, error) {
	return b.bus.WaitFor[T](ctx, match)
}

// Stream yields packets or events of type T until ctx is done or the loop stops.
func (b *Client) Stream[T any](ctx context.Context) iter.Seq[T] {
	return b.bus.Stream[T](ctx)
}

// Emit publishes an event to its subscribers.
func (b *Client) Emit[T any](ctx context.Context, v T) {
	b.bus.Emit[T](ctx, v)
}

// OnFrame observes every incoming packet before decoding.
func (b *Client) OnFrame(f bus.FrameFunc) (off func()) {
	return b.bus.OnFrame(f)
}

// OnRaw subscribes to the undecoded payload of one packet ID. p.Data is only
// valid until f returns.
func (b *Client) OnRaw(id packetid.ClientboundPacketID, f func(ctx context.Context, p pk.Packet)) (off func()) {
	return b.bus.OnRaw(id, f)
}

func (b *Client) World() bot.World {
	return b.world
}

func (b *Client) Inventory() bot.InventoryHandler {
	return b.inventory
}

func (b *Client) Connect(ctx context.Context, addr string, options *bot.ConnectOptions) error {
	b.connMu.RLock()
	hasConnection := b.conn != nil
	b.connMu.RUnlock()
	if hasConnection {
		return errors.New("client already has an open connection")
	}
	b.registryMu.Lock()
	b.registryEntries = nil
	b.registryMu.Unlock()
	b.setConnectionState(packet.StateLogin)

	host, portStr, err := net.SplitHostPort(addr)
	var port uint64
	if err != nil {
		var addrErr *net.AddrError
		const missingPort = "missing port in address"
		if errors.As(err, &addrErr) && addrErr.Err == missingPort {
			host = addr
			port = 25565
		} else {
			return err
		}
	} else {
		port, err = strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return err
		}
	}

	var dialer mcnet.MCDialer = &mcnet.DefaultDialer
	if options != nil && options.Proxy != nil {
		dialer, err = socks5(options.Proxy)
		if err != nil {
			return err
		}
	}
	conn, err := dialer.DialMCContext(ctx, addr)
	if err != nil {
		return err
	}
	b.connMu.Lock()
	b.conn = conn
	b.connMu.Unlock()
	connected := false
	defer func() {
		if !connected {
			_ = b.Close(context.Background())
		}
	}()

	if options != nil && options.FakeHost != "" {
		host = options.FakeHost
	}

	err = b.handshake(ctx, host, port)
	if err != nil {
		return err
	}

	err = b.login(ctx)
	if err != nil {
		return err
	}

	b.setConnectionState(packet.StateConfig)
	b.bus.Emit(ctx, ConnectionStateChangeEvent{From: packet.StateLogin, To: packet.StateConfig})

	err = b.configuration(ctx)
	if err != nil {
		return err
	}

	b.setConnectionState(packet.StatePlay)
	b.bus.Emit(ctx, ConnectionStateChangeEvent{From: packet.StateConfig, To: packet.StatePlay})

	b.connected.Store(true)
	connected = true

	return nil
}

func (b *Client) HandleGame(ctx context.Context) error {
	defer func() {
		b.connected.Store(false)
		b.registryMu.Lock()
		b.registryEntries = nil
		b.registryMu.Unlock()
	}()
	return b.handlePackets(ctx)
}

func (b *Client) handshake(ctx context.Context, host string, port uint64) error {
	return b.writeRawPacket(ctx, pk.Marshal(
		0,
		pk.VarInt(protocol.ProtocolVersion),
		pk.String(host),
		pk.UnsignedShort(port),
		pk.VarInt(2), // to game state
	))
}

func (b *Client) handlePackets(ctx context.Context) error {
	b.connMu.RLock()
	conn := b.conn
	b.connMu.RUnlock()
	if conn == nil {
		return errors.New("client is not connected")
	}

	stopRead := context.AfterFunc(ctx, func() {
		_ = conn.Socket.SetReadDeadline(time.Now())
	})
	defer func() {
		stopRead()
		_ = conn.Socket.SetReadDeadline(time.Time{})
	}()

	// p is reused for every packet: handlers run synchronously and decoded
	// packets copy what they keep, so nothing refers to p.Data after dispatch.
	var p pk.Packet
	var reader bytes.Reader
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
		}

		// Refresh the idle deadline only when the next read may block on the
		// socket; while packets are buffered it costs a clock read and a
		// timer reset per packet for nothing.
		if conn.Buffered() == 0 {
			if err := b.setReadDeadline(ctx, conn); err != nil {
				return err
			}
		}
		if err := conn.ReadPacket(&p); err != nil {
			if ctx.Err() != nil {
				return context.Cause(ctx)
			}
			return err
		}
		pktID := packetid.ClientboundPacketID(p.ID)
		if pktID == packetid.ClientboundStartConfiguration {
			if err := b.reconfigure(ctx); err != nil {
				return err
			}
			continue
		}

		b.bus.DispatchFrame(pktID, len(p.Data))
		b.bus.DispatchRaw(ctx, pktID, p)

		// Decode only what someone subscribed to.
		if !b.bus.WantsPacket(pktID) {
			continue
		}
		creator, ok := client.ClientboundPackets[pktID]
		if !ok {
			continue
		}
		pkt := creator()
		reader.Reset(p.Data)
		if err := decodePacket(pkt, &reader); err != nil {
			return b.decodeError(pktID, p.Data, len(p.Data)-reader.Len(), err)
		}
		b.bus.DispatchPacket(ctx, pkt)
	}
}

// reconfigure handles the server sending the client back to configuration.
func (b *Client) reconfigure(ctx context.Context) error {
	b.setConnectionState(packet.StateConfig)
	b.registryMu.Lock()
	b.registryEntries = nil
	b.registryMu.Unlock()
	b.bus.Emit(ctx, ConnectionStateChangeEvent{From: packet.StatePlay, To: packet.StateConfig})

	if err := b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigurationAcknowledged)); err != nil {
		return err
	}
	if err := b.configuration(ctx); err != nil {
		return err
	}

	b.setConnectionState(packet.StatePlay)
	b.bus.Emit(ctx, ConnectionStateChangeEvent{From: packet.StateConfig, To: packet.StatePlay})
	return nil
}

// decodePacket decodes one packet, turning a decoder panic into an error so
// one malformed packet cannot take the whole process down.
func decodePacket(pkt client.ClientboundPacket, r *bytes.Reader) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("decoder panicked: %v", rec)
		}
	}()
	_, err = pkt.ReadFrom(r)
	return err
}

// decodeError dumps the packet that failed to decode to a temp file for
// debugging and returns the error that ends the connection.
func (b *Client) decodeError(pktID packetid.ClientboundPacketID, data []byte, at int, err error) error {
	temp, dumpErr := os.CreateTemp("", "packet")
	filename := ""
	if dumpErr == nil {
		filename = temp.Name()
		_, dumpErr = temp.Write(data)
		if closeErr := temp.Close(); dumpErr == nil {
			dumpErr = closeErr
		}
	}
	slog.Error("decode clientbound packet", "packet", pktID, "dataLen", len(data), "err", err, "filename", filename, "dumpErr", dumpErr)
	return fmt.Errorf("decode clientbound packet %d at byte %d/%d: %w TEMPPACKET: %s", pktID, at, len(data), err, filename)
}

// NewClient returns a client that is not connected yet. The returned
// *Client is also a bot.Client.
func NewClient(options *bot.ClientOptions) *Client {
	c := &Client{
		bus:                bus.New(),
		readIdleTimeout:    30 * time.Second,
		resourcePackPolicy: bot.ResourcePackAccept,
		cookies:            make(map[string][]int8),
		state:              packet.StatePlay,
		stateChanged:       make(chan struct{}),
	}

	if options != nil {
		c.authProvider = options.AuthProvider
		if options.ReadIdleTimeout > 0 {
			c.readIdleTimeout = options.ReadIdleTimeout
		}
		if options.ResourcePackPolicy != "" {
			c.resourcePackPolicy = options.ResourcePackPolicy
		}
	}
	if c.authProvider == nil {
		c.authProvider = &auth.OfflineAuth{Username: "Steve"}
	}

	var worldOpts []world.Option
	if options != nil && options.NoChunks {
		worldOpts = append(worldOpts, world.WithoutChunks())
	}
	c.world = world.NewWorld(c, worldOpts...)
	c.inventory = inventory.NewManager(c)
	var playerOpts []player.Option
	if options != nil && options.NoTeleport {
		playerOpts = append(playerOpts, player.WithoutTeleport())
	}
	c.player = player.New(c, playerOpts...)
	c.On(func(ctx context.Context, packet *client.AddResourcePack) {
		for _, result := range resourcePackResults(c.resourcePackPolicy) {
			if err := c.WritePacket(ctx, &server.ResourcePack{UUID: packet.UUID, Result: result}); err != nil {
				slog.Error("respond to play resource pack", "uuid", packet.UUID, "result", result, "err", err)
				return
			}
		}
	})
	c.On(func(_ context.Context, packet *client.StoreCookie) {
		c.storeCookie(packet.Key, packet.Payload)
	})
	c.On(func(ctx context.Context, packet *client.CookieRequest) {
		response := c.cookieResponse(string(packet.Key))
		if err := c.WritePacket(ctx, &response); err != nil {
			slog.Error("respond to play cookie request", "key", packet.Key, "err", err)
		}
	})

	return c
}

func resourcePackResults(policy bot.ResourcePackPolicy) []int32 {
	if policy == bot.ResourcePackDecline {
		return []int32{1}
	}
	// 3 = accepted, 0 = successfully loaded. Minego is headless and does not
	// render assets, so it can safely acknowledge without downloading them.
	return []int32{3, 0}
}

func (b *Client) storeCookie(key string, payload []int8) {
	b.cookieMu.Lock()
	defer b.cookieMu.Unlock()
	if b.cookies == nil {
		b.cookies = make(map[string][]int8)
	}
	b.cookies[key] = append([]int8(nil), payload...)
}

func (b *Client) cookieResponse(key string) server.CookieResponse {
	b.cookieMu.RLock()
	payload, ok := b.cookies[key]
	b.cookieMu.RUnlock()
	return server.CookieResponse{
		Key:        key,
		HasPayload: ok,
		Payload:    append([]int8(nil), payload...),
	}
}
