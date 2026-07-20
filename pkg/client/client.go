package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"
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
	conn          *mcnet.Conn
	connMu        sync.RWMutex
	writeMu       sync.Mutex
	packetHandler *packetHandler
	eventHandler  bot.EventHandler
	connected     atomic.Bool
	authProvider  auth.Provider

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

func (b *botClient) WritePacket(ctx context.Context, packet server.ServerboundPacket) error {
	return b.writeRawPacket(ctx, pk.Marshal(packet.PacketID(), packet))
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

	err = b.eventHandler.PublishEvent(EventConnectionStateChange, ConnectionStateChangeEvent{From: packet.StateLogin, To: packet.StateConfig})
	if err != nil {
		return err
	}

	err = b.configuration(ctx)
	if err != nil {
		return err
	}

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
		pk.VarInt(776), // TODO 版本更新時要記得改 current: 26.2
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

	const readTimeout = 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			var p pk.Packet

			readDeadline := time.Now().Add(readTimeout)
			if deadline, ok := ctx.Deadline(); ok && deadline.Before(readDeadline) {
				readDeadline = deadline
			}
			if err := conn.Socket.SetReadDeadline(readDeadline); err != nil {
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
				handlers.Go(func() {
					defer func() { <-semaphore }()
					h(handlerCtx, p)
				})
			}

			if _, ok := client.ClientboundPackets[pktID]; !ok {
				continue
			}
			pkt, handled, err := decodeClientboundPacket(pktID, p.Data)
			if err != nil {
				return fmt.Errorf("decode clientbound packet %d: %w", pktID, err)
			}
			if !handled {
				continue
			}
			b.packetHandler.handlePacket(ctx, pktID, pkt)

			_ = conn.Socket.SetReadDeadline(time.Time{})
		}
	}
}

func decodeClientboundPacket(id packetid.ClientboundPacketID, data []byte) (client.ClientboundPacket, bool, error) {
	creator, ok := client.ClientboundPackets[id]
	if !ok {
		return nil, false, nil
	}

	pkt := creator()
	reader := bytes.NewReader(data)
	_, err := pkt.ReadFrom(reader)
	if err != nil {
		return nil, true, err
	}
	if reader.Len() != 0 {
		return nil, true, fmt.Errorf("decoder left %d of %d bytes unread", reader.Len(), len(data))
	}
	return pkt, true, nil
}

func NewClient(options *bot.ClientOptions) bot.Client {
	c := &botClient{
		packetHandler: newPacketHandler(),
		eventHandler:  NewEventHandler(),
	}

	if options != nil {
		c.authProvider = options.AuthProvider
	}
	if c.authProvider == nil {
		c.authProvider = &auth.OfflineAuth{Username: "Steve"}
	}

	c.world = world.NewWorld(c)
	c.inventory = inventory.NewManager(c)
	c.player = player.New(c)

	return c
}
