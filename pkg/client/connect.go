package client

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/auth"
)

const resourcePackResultDeclined = 1

func (b *botClient) login(ctx context.Context) error {
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	return b.withReadContext(ctx, func(conn *mcnet.Conn) error {
		a := &auth.Auth{Conn: conn, Provider: b.authProvider}
		return a.HandleLogin(ctx)
	})
}

func (b *botClient) configuration(ctx context.Context) (err error) {
	return b.withReadContext(ctx, func(conn *mcnet.Conn) error {
		return b.readConfiguration(ctx, conn)
	})
}

func (b *botClient) readConfiguration(ctx context.Context, conn *mcnet.Conn) (err error) {
	var p pk.Packet
	for {
		err = conn.ReadPacket(&p)
		if err != nil {
			return err
		}

		switch packetid.ClientboundPacketID(p.ID) {
		case packetid.ClientboundConfigDisconnect:
			var reason chat.Message
			err = p.Scan(&reason)
			if err != nil {
				return err
			}
			return errors.New("kicked: " + reason.String())
		case packetid.ClientboundConfigFinishConfiguration:
			err = b.writeRawPacket(ctx, pk.Marshal(
				packetid.ServerboundConfigFinishConfiguration,
			))
			return err
		case packetid.ClientboundConfigKeepAlive:
			var keepAliveID pk.Long
			err = p.Scan(&keepAliveID)
			if err != nil {
				return err
			}
			err = b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigKeepAlive, keepAliveID))
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigPing:
			var pingID pk.Int
			err = p.Scan(&pingID)
			if err != nil {
				return err
			}
			err = b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigPong, pingID))
			if err != nil {
				return err
			}

		case packetid.ClientboundConfigResourcePackPush:
			var packID pk.UUID
			err = p.Scan(&packID)
			if err != nil {
				return err
			}
			err = b.writeRawPacket(ctx, pk.Marshal(
				packetid.ServerboundConfigResourcePack,
				packID,
				pk.VarInt(resourcePackResultDeclined),
			))
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigSelectKnownPacks:
			err = b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigSelectKnownPacks, pk.VarInt(0)))
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigCodeOfConduct:
			// Consume the packet payload so malformed packets still fail cleanly.
			var codeOfConduct pk.String
			if _, err = (&codeOfConduct).ReadFrom(bytes.NewReader(p.Data)); err != nil {
				return err
			}
			err = b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigAcceptCodeOfConduct))
			if err != nil {
				return err
			}
		default:
			continue
		}
	}
}

func (b *botClient) withReadContext(ctx context.Context, fn func(*mcnet.Conn) error) error {
	b.connMu.RLock()
	conn := b.conn
	b.connMu.RUnlock()
	if conn == nil {
		return errors.New("client is not connected")
	}

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.Socket.SetReadDeadline(deadline); err != nil {
			return err
		}
	}
	stop := context.AfterFunc(ctx, func() {
		_ = conn.Socket.SetReadDeadline(time.Now())
	})
	defer func() {
		stop()
		_ = conn.Socket.SetReadDeadline(time.Time{})
	}()

	err := fn(conn)
	if ctx.Err() != nil {
		return context.Cause(ctx)
	}
	return err
}
