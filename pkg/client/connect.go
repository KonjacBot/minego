package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/auth"
	"github.com/KonjacBot/minego/pkg/protocol"
	configclient "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/client"
	configserver "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/server"
	gameserver "github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

func (b *botClient) login(ctx context.Context) error {
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	return b.withReadContext(ctx, func(conn *mcnet.Conn) error {
		a := &auth.Auth{
			Conn:            conn,
			Provider:        b.authProvider,
			ReadIdleTimeout: b.readIdleTimeout,
		}
		return a.HandleLogin(ctx)
	})
}

func (b *botClient) configuration(ctx context.Context) (err error) {
	return b.withReadContext(ctx, func(conn *mcnet.Conn) error {
		return b.readConfiguration(ctx, conn)
	})
}

func (b *botClient) readConfiguration(ctx context.Context, conn *mcnet.Conn) (err error) {
	information := configserver.ConfigClientInformation{
		ClientInformation: gameserver.ClientInformation{
			Location:            "zh_TW",
			ViewDistance:        16,
			ChatMode:            0,
			ChatColor:           true,
			DisplayedSkinParts:  127,
			MainHand:            0,
			EnableTextFiltering: false,
			AllowListing:        true,
			ParticleStatus:      0,
		},
	}
	if err = b.writeRawPacket(ctx, pk.Marshal(information.PacketID(), &information)); err != nil {
		return fmt.Errorf("write configuration client information: %w", err)
	}

	var p pk.Packet
	for {
		if err = b.setReadDeadline(ctx, conn); err != nil {
			return err
		}
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

		case packetid.ClientboundConfigSelectKnownPacks:
			var offered configclient.ConfigSelectKnownPacks
			err = p.Scan(&offered)
			if err != nil {
				return err
			}
			response := configserver.ConfigSelectKnownPacks{
				Packs: supportedKnownPacks(offered.KnownPacks),
			}
			err = b.writeRawPacket(ctx, pk.Marshal(packetid.ServerboundConfigSelectKnownPacks, &response))
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigStoreCookie:
			var stored configclient.ConfigStoreCookie
			if err = p.Scan(&stored); err != nil {
				return err
			}
			b.storeCookie(stored.Key, stored.Payload)
		case packetid.ClientboundConfigCookieRequest:
			var request configclient.ConfigCookieRequest
			if err = p.Scan(&request); err != nil {
				return err
			}
			gameResponse := b.cookieResponse(request.Key)
			response := configserver.ConfigCookieResponse{CookieResponse: gameResponse}
			if err = b.writeRawPacket(ctx, pk.Marshal(response.PacketID(), &response)); err != nil {
				return err
			}
		case packetid.ClientboundConfigResourcePackPush:
			var offered configclient.ConfigResourcePackPush
			if err = p.Scan(&offered); err != nil {
				return err
			}
			for _, result := range resourcePackResults(b.resourcePackPolicy) {
				response := configserver.ConfigResourcePack{
					ResourcePack: gameserver.ResourcePack{
						UUID:   offered.UUID,
						Result: result,
					},
				}
				if err = b.writeRawPacket(ctx, pk.Marshal(response.PacketID(), &response)); err != nil {
					return err
				}
			}
		default:
			continue
		}
	}
}

func supportedKnownPacks(offered []configclient.KnownPack) []configclient.KnownPack {
	result := make([]configclient.KnownPack, 0, 1)
	for _, pack := range offered {
		if pack.Namespace == protocol.CorePackNamespace &&
			pack.ID == protocol.CorePackID &&
			pack.Version == protocol.VersionName {
			result = append(result, pack)
		}
	}
	return result
}

func (b *botClient) withReadContext(ctx context.Context, fn func(*mcnet.Conn) error) error {
	b.connMu.RLock()
	conn := b.conn
	b.connMu.RUnlock()
	if conn == nil {
		return errors.New("client is not connected")
	}

	if err := b.setReadDeadline(ctx, conn); err != nil {
		return err
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

func (b *botClient) setReadDeadline(ctx context.Context, conn *mcnet.Conn) error {
	var deadline time.Time
	if b.readIdleTimeout > 0 {
		deadline = time.Now().Add(b.readIdleTimeout)
	}
	if contextDeadline, ok := ctx.Deadline(); ok && (deadline.IsZero() || contextDeadline.Before(deadline)) {
		deadline = contextDeadline
	}
	return conn.Socket.SetReadDeadline(deadline)
}
