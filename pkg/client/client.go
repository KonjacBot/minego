package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/KonjacBot/minego/pkg/protocol/packet"

	"github.com/KonjacBot/minego/pkg/auth"
	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/game/inventory"
	"github.com/KonjacBot/minego/pkg/game/player"
	"github.com/KonjacBot/minego/pkg/game/world"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

type botClient struct {
	conn               *mcnet.Conn
	connMu             sync.RWMutex
	writeMu            sync.Mutex
	stateMu            sync.Mutex
	state              packet.State
	stateChanged       chan struct{}
	packetHandler      *packetHandler
	eventHandler       bot.EventHandler
	connected          atomic.Bool
	authProvider       auth.Provider
	readIdleTimeout    time.Duration
	resourcePackPolicy bot.ResourcePackPolicy
	cookieMu           sync.RWMutex
	cookies            map[string][]int8

	inventory *inventory.Manager
	world     *world.World
	player    *player.Player
}

func (b *botClient) Player() bot.Player {
	return b.player
}

func (b *botClient) Close(ctx context.Context) error {
	b.writeMu.Lock()
	defer b.writeMu.Unlock()

	b.connMu.Lock()
	conn := b.conn
	b.conn = nil
	b.connected.Store(false)
	b.connMu.Unlock()
	b.setConnectionStateLocked(packet.StateLogin)

	if conn == nil {
		return ctx.Err()
	}
	if err := conn.Close(); err != nil {
		return err
	}
	return ctx.Err()
}

func (b *botClient) IsConnected() bool {
	return b.connected.Load()
}

func (b *botClient) WritePacket(ctx context.Context, outbound server.ServerboundPacket) error {
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

func (b *botClient) writeRawPacket(ctx context.Context, packet pk.Packet) error {
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

func (b *botClient) writeRawPacketLocked(ctx context.Context, packet pk.Packet) error {
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

func (b *botClient) waitForPlayState(ctx context.Context) error {
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

func (b *botClient) connectionState() packet.State {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	return b.state
}

func (b *botClient) setConnectionState(state packet.State) {
	b.writeMu.Lock()
	defer b.writeMu.Unlock()
	b.setConnectionStateLocked(state)
}

func (b *botClient) setConnectionStateLocked(state packet.State) {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	b.state = state
	if b.stateChanged != nil {
		close(b.stateChanged)
	}
	b.stateChanged = make(chan struct{})
}

func (b *botClient) PacketHandler() bot.PacketHandler {
	return b.packetHandler
}

func (b *botClient) EventHandler() bot.EventHandler {
	return b.eventHandler
}

func (b *botClient) World() bot.World {
	return b.world
}

func (b *botClient) Inventory() bot.InventoryHandler {
	return b.inventory
}

func (b *botClient) Connect(ctx context.Context, addr string, options *bot.ConnectOptions) error {
	b.connMu.RLock()
	hasConnection := b.conn != nil
	b.connMu.RUnlock()
	if hasConnection {
		return errors.New("client already has an open connection")
	}
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
	err = b.eventHandler.PublishEvent(EventConnectionStateChange, ConnectionStateChangeEvent{From: packet.StateLogin, To: packet.StateConfig})
	if err != nil {
		return err
	}

	err = b.configuration(ctx)
	if err != nil {
		return err
	}

	b.setConnectionState(packet.StatePlay)
	err = b.eventHandler.PublishEvent(EventConnectionStateChange, ConnectionStateChangeEvent{From: packet.StateConfig, To: packet.StatePlay})
	if err != nil {
		return err
	}

	b.connected.Store(true)
	connected = true

	return nil
}

func (b *botClient) HandleGame(ctx context.Context) error {
	defer b.connected.Store(false)
	return b.handlePackets(ctx)
}

func (b *botClient) handshake(ctx context.Context, host string, port uint64) error {
	return b.writeRawPacket(ctx, pk.Marshal(
		0,
		pk.VarInt(protocol.ProtocolVersion),
		pk.String(host),
		pk.UnsignedShort(port),
		pk.VarInt(2), // to game state
	))
}

func (b *botClient) handlePackets(ctx context.Context) error {
	b.connMu.RLock()
	conn := b.conn
	b.connMu.RUnlock()
	if conn == nil {
		return errors.New("client is not connected")
	}

	handlerCtx, cancelHandlers := context.WithCancel(ctx)
	var handlers sync.WaitGroup
	semaphore := make(chan struct{}, 15)
	defer func() {
		cancelHandlers()
		handlers.Wait()
	}()
	stopRead := context.AfterFunc(ctx, func() {
		_ = conn.Socket.SetReadDeadline(time.Now())
	})
	defer func() {
		stopRead()
		_ = conn.Socket.SetReadDeadline(time.Time{})
	}()

	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			var p pk.Packet

			if err := b.setReadDeadline(ctx, conn); err != nil {
				return err
			}

			if err := conn.ReadPacket(&p); err != nil {
				if ctx.Err() != nil {
					return context.Cause(ctx)
				}
				return err
			}
			pktID := packetid.ClientboundPacketID(p.ID)
			if pktID == packetid.ClientboundStartConfiguration {
				b.setConnectionState(packet.StateConfig)
				err := b.eventHandler.PublishEvent(EventConnectionStateChange, ConnectionStateChangeEvent{From: packet.StatePlay, To: packet.StateConfig})
				if err != nil {
					return err
				}

				err = b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigurationAcknowledged))
				if err != nil {
					return err
				}

				err = b.configuration(ctx)
				if err != nil {
					return err
				}

				b.setConnectionState(packet.StatePlay)
				err = b.eventHandler.PublishEvent(EventConnectionStateChange, ConnectionStateChangeEvent{From: packet.StateConfig, To: packet.StatePlay})
				if err != nil {
					return err
				}
				continue
			}

			hs := b.packetHandler.rawHandlers(pktID)
			for _, h := range hs {
				select {
				case semaphore <- struct{}{}:
				case <-ctx.Done():
					return context.Cause(ctx)
				}
				packetID := pktID
				packetCopy := p
				handler := h
				handlers.Go(func() {
					defer func() { <-semaphore }()
					defer func() {
						if recovered := recover(); recovered != nil {
							slog.Error("raw packet handler panic",
								"packet", packetID,
								"error", recovered,
								"stack", string(debug.Stack()),
							)
						}
					}()
					handler(handlerCtx, packetCopy)
				})
			}

			creator, ok := client.ClientboundPackets[pktID]
			if !ok {
				continue
			}
			pkt := creator()
			reader := bytes.NewReader(p.Data)
			_, err := pkt.ReadFrom(reader)
			if err != nil {
				temp, dumpErr := os.CreateTemp("", "packet")
				filename := ""
				if dumpErr == nil {
					filename = temp.Name()
					_, dumpErr = temp.Write(p.Data)
					if closeErr := temp.Close(); dumpErr == nil {
						dumpErr = closeErr
					}
				}
				slog.Error("decode clientbound packet", "packet", pktID, "dataLen", len(p.Data), "err", err, "filename", filename, "dumpErr", dumpErr)
				return fmt.Errorf("decode clientbound packet %d at byte %d/%d: %w TEMPPACKET: %s", pktID, len(p.Data)-reader.Len(), len(p.Data), err, filename)
			}
			b.packetHandler.HandlePacket(ctx, pkt)

			_ = conn.Socket.SetReadDeadline(time.Time{})
		}
	}
}

func NewClient(options *bot.ClientOptions) bot.Client {
	c := &botClient{
		packetHandler:      newPacketHandler(),
		eventHandler:       NewEventHandler(),
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

	c.world = world.NewWorld(c)
	c.inventory = inventory.NewManager(c)
	c.player = player.New(c)
	bot.AddHandler(c, func(ctx context.Context, packet *client.AddResourcePack) {
		for _, result := range resourcePackResults(c.resourcePackPolicy) {
			if err := c.WritePacket(ctx, &server.ResourcePack{UUID: packet.UUID, Result: result}); err != nil {
				slog.Error("respond to play resource pack", "uuid", packet.UUID, "result", result, "err", err)
				return
			}
		}
	})
	bot.AddHandler(c, func(_ context.Context, packet *client.StoreCookie) {
		c.storeCookie(packet.Key, packet.Payload)
	})
	bot.AddHandler(c, func(ctx context.Context, packet *client.CookieRequest) {
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

func (b *botClient) storeCookie(key string, payload []int8) {
	b.cookieMu.Lock()
	defer b.cookieMu.Unlock()
	if b.cookies == nil {
		b.cookies = make(map[string][]int8)
	}
	b.cookies[key] = append([]int8(nil), payload...)
}

func (b *botClient) cookieResponse(key string) server.CookieResponse {
	b.cookieMu.RLock()
	payload, ok := b.cookies[key]
	b.cookieMu.RUnlock()
	return server.CookieResponse{
		Key:        key,
		HasPayload: ok,
		Payload:    append([]int8(nil), payload...),
	}
}
