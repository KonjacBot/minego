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
	conn                 *mcnet.Conn
	connMu               sync.RWMutex
	connectMu            sync.Mutex
	writeMu              sync.Mutex
	packetHandler        *packetHandler
	eventHandler         bot.EventHandler
	connected            atomic.Bool
	handling             atomic.Bool
	compressionSet       atomic.Bool
	compressionThreshold atomic.Int32
	authProvider         auth.Provider

	inventory *inventory.Manager
	world     *world.World
	player    *player.Player
}

func (b *botClient) Player() bot.Player {
	return b.player
}

func (b *botClient) Close(ctx context.Context) error {
	b.connMu.Lock()
	conn := b.conn
	b.conn = nil
	b.connected.Store(false)
	b.connMu.Unlock()

	if conn == nil {
		return ctx.Err()
	}
	return conn.Close()
}

func (b *botClient) IsConnected() bool {
	return b.connected.Load()
}

func (b *botClient) WritePacket(ctx context.Context, packet server.ServerboundPacket) error {
	if err := ctx.Err(); err != nil {
		return context.Cause(ctx)
	}
	encoded, err := marshalServerboundPacket(packet)
	if err != nil {
		return err
	}
	return b.writeRawPacket(ctx, encoded)
}

type packetValidator interface {
	Validate() error
}

func marshalServerboundPacket(value server.ServerboundPacket) (result pk.Packet, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = pk.Packet{}
			err = fmt.Errorf("encode serverbound packet %T panicked: %v", value, recovered)
		}
	}()

	if value == nil {
		return pk.Packet{}, errors.New("serverbound packet is nil")
	}
	if validator, ok := value.(packetValidator); ok {
		if err := validator.Validate(); err != nil {
			return pk.Packet{}, fmt.Errorf("validate serverbound packet %T: %w", value, err)
		}
	}
	return packet.Marshal(value.PacketID(), value)
}

func (b *botClient) writePacketFields(ctx context.Context, id int32, fields ...pk.FieldEncoder) error {
	encoded, err := packet.Marshal(id, fields...)
	if err != nil {
		return err
	}
	return b.writeRawPacket(ctx, encoded)
}

func (b *botClient) writeRawPacket(ctx context.Context, value pk.Packet) error {
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
	}
	writeCanceled := make(chan struct{})
	stopCancel := context.AfterFunc(ctx, func() {
		_ = conn.Socket.SetWriteDeadline(time.Now())
		close(writeCanceled)
	})
	defer func() {
		if !stopCancel() {
			<-writeCanceled
		}
		_ = conn.Socket.SetWriteDeadline(time.Time{})
	}()

	err := packet.WriteFrame(conn.Writer, b.frameThreshold(), value)
	if ctx.Err() != nil {
		return context.Cause(ctx)
	}
	return err
}

func (b *botClient) frameThreshold() int {
	if !b.compressionSet.Load() {
		return -1
	}
	return int(b.compressionThreshold.Load())
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
	b.connectMu.Lock()
	defer b.connectMu.Unlock()

	b.connMu.RLock()
	hasConnection := b.conn != nil
	b.connMu.RUnlock()
	if hasConnection {
		return errors.New("client already has an open connection")
	}
	b.compressionSet.Store(false)

	host, portStr, err := net.SplitHostPort(addr)
	var port uint64
	if err != nil {
		var addrErr *net.AddrError
		const missingPort = "missing port in address"
		if (errors.As(err, &addrErr) && addrErr.Err == missingPort) || net.ParseIP(addr) != nil {
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
	dialAddress := addr
	if options != nil && options.Proxy != nil {
		dialer, err = socks5(options.Proxy)
		if err != nil {
			return err
		}
		dialAddress = net.JoinHostPort(host, strconv.FormatUint(port, 10))
	}
	conn, err := dialer.DialMCContext(ctx, dialAddress)
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
	if !b.handling.CompareAndSwap(false, true) {
		return errors.New("client is already handling game packets")
	}
	defer b.handling.Store(false)
	defer b.connected.Store(false)
	return b.handlePackets(ctx)
}

func (b *botClient) handshake(ctx context.Context, host string, port uint64) error {
	return b.writePacketFields(ctx, 0,
		pk.VarInt(776), // TODO 版本更新時要記得改 current: 26.2
		pk.String(host),
		pk.UnsignedShort(port),
		pk.VarInt(2), // to game state
	)
}

func (b *botClient) handlePackets(ctx context.Context) error {
	b.connMu.RLock()
	conn := b.conn
	b.connMu.RUnlock()
	if conn == nil {
		return errors.New("client is not connected")
	}

	handlerCtx, cancelHandlers := context.WithCancel(ctx)
	semaphore := make(chan struct{}, 15)
	handlerErr := make(chan error, 1)
	defer cancelHandlers()
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
		case err := <-handlerErr:
			return err
		default:
			var p pk.Packet

			readDeadline := time.Now().Add(readTimeout)
			if deadline, ok := ctx.Deadline(); ok && deadline.Before(readDeadline) {
				readDeadline = deadline
			}
			if err := conn.Socket.SetReadDeadline(readDeadline); err != nil {
				return err
			}

			if err := packet.ReadFrame(conn.Reader, b.frameThreshold(), &p); err != nil {
				select {
				case handlerFailure := <-handlerErr:
					return handlerFailure
				default:
				}
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

				err = b.writePacketFields(ctx, int32(packetid.ServerboundConfigurationAcknowledged))
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

			pkt, handled, err := decodeClientboundPacket(pktID, p.Data)
			if err != nil {
				return fmt.Errorf("decode clientbound packet %d: %w", pktID, err)
			}

			hs := b.packetHandler.rawHandlers(pktID)
			for _, h := range hs {
				select {
				case semaphore <- struct{}{}:
				case <-ctx.Done():
					return context.Cause(ctx)
				}
				rawPacket := pk.Packet{ID: p.ID, Data: bytes.Clone(p.Data)}
				go func(handler func(context.Context, pk.Packet), packet pk.Packet) {
					defer func() { <-semaphore }()
					if callbackErr := invokeCallback("raw packet handler", func() {
						handler(handlerCtx, packet)
					}); callbackErr != nil {
						select {
						case handlerErr <- callbackErr:
							_ = conn.Socket.SetReadDeadline(time.Now())
						default:
						}
					}
				}(h, rawPacket)
			}

			if !handled {
				continue
			}
			if err := b.packetHandler.handlePacket(ctx, pktID, pkt); err != nil {
				return err
			}
			if err := b.handleControlPacket(ctx, pkt); err != nil {
				return err
			}

			_ = conn.Socket.SetReadDeadline(time.Time{})
		}
	}
}

type DisconnectError struct {
	Reason string
}

func (e *DisconnectError) Error() string {
	return "server disconnected: " + e.Reason
}

func (b *botClient) handleControlPacket(ctx context.Context, value client.ClientboundPacket) error {
	switch p := value.(type) {
	case *client.Disconnect:
		return &DisconnectError{Reason: p.Reason.String()}
	case *client.CookieRequest:
		return b.WritePacket(ctx, &server.CookieResponse{Key: string(p.Key)})
	case *client.ChunkBatchFinished:
		return b.WritePacket(ctx, &server.ChunkBatchReceived{ChunksPerTick: 7})
	case *client.AddResourcePack:
		return b.WritePacket(ctx, &server.ResourcePack{UUID: p.UUID, Result: resourcePackResultDeclined})
	default:
		return nil
	}
}

func decodeClientboundPacket(id packetid.ClientboundPacketID, data []byte) (pkt client.ClientboundPacket, handled bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			pkt = nil
			handled = true
			err = fmt.Errorf("decoder panicked: %v", recovered)
		}
	}()

	pkt, handled = client.NewClientboundPacket(id)
	if !handled {
		return nil, false, nil
	}
	if pkt == nil {
		return nil, true, errors.New("packet registry returned a nil packet")
	}
	if err = decodeClientboundPayload(pkt, data); err != nil {
		return nil, true, err
	}
	return pkt, true, nil
}

func decodeClientboundPayload(pkt client.ClientboundPacket, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := pkt.ReadFrom(reader)
	if err != nil {
		return err
	}
	if reader.Len() != 0 {
		return fmt.Errorf("decoder left %d of %d bytes unread", reader.Len(), len(data))
	}
	return nil
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
