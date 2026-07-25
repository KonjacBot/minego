package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/auth"
	protocolpacket "github.com/KonjacBot/minego/pkg/protocol/packet"
	configclient "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/client"
	configserver "github.com/KonjacBot/minego/pkg/protocol/packet/configuration/server"
	gameserver "github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

const resourcePackResultDeclined = 1

func (b *botClient) login(ctx context.Context) error {
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	return b.withReadContext(ctx, func(conn *mcnet.Conn) error {
		a := &auth.Auth{
			Conn:     conn,
			Provider: b.authProvider,
			OnSetCompression: func(threshold int) {
				b.compressionThreshold.Store(int32(threshold))
				b.compressionSet.Store(true)
			},
		}
		return a.HandleLogin(ctx)
	})
}

func (b *botClient) configuration(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return b.withReadContext(ctx, func(conn *mcnet.Conn) error {
		return b.readConfiguration(ctx, conn)
	})
}

func (b *botClient) readConfiguration(ctx context.Context, conn *mcnet.Conn) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("configuration callback or decoder panicked: %v", recovered)
		}
	}()

	err = b.WritePacket(ctx, &configserver.ConfigClientInformation{
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
	})
	if err != nil {
		return fmt.Errorf("write configuration client information: %w", err)
	}

	var p pk.Packet
	for {
		err = protocolpacket.ReadFrame(conn.Reader, b.frameThreshold(), &p)
		if err != nil {
			return err
		}

		switch packetid.ClientboundPacketID(p.ID) {
		case packetid.ClientboundConfigDisconnect:
			var disconnect configclient.ConfigDisconnect
			err = protocolpacket.Scan(p, &disconnect)
			if err != nil {
				return err
			}
			return errors.New("kicked: " + disconnect.Reason.String())
		case packetid.ClientboundConfigFinishConfiguration:
			if err = protocolpacket.Scan(p); err != nil {
				return err
			}
			err = b.WritePacket(ctx, &configserver.ConfigFinishConfiguration{})
			return err
		case packetid.ClientboundConfigKeepAlive:
			var keepAlive configclient.ConfigKeepAlive
			err = protocolpacket.Scan(p, &keepAlive)
			if err != nil {
				return err
			}
			err = b.WritePacket(ctx, &configserver.ConfigKeepAlive{ID: keepAlive.ID})
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigPing:
			var ping configclient.ConfigPing
			err = protocolpacket.Scan(p, &ping)
			if err != nil {
				return err
			}
			err = b.WritePacket(ctx, &configserver.ConfigPong{ID: ping.ID})
			if err != nil {
				return err
			}

		case packetid.ClientboundConfigResourcePackPush:
			var resourcePack configclient.ConfigResourcePackPush
			err = protocolpacket.Scan(p, &resourcePack)
			if err != nil {
				return err
			}
			err = b.WritePacket(ctx, &configserver.ConfigResourcePack{
				ResourcePack: gameserver.ResourcePack{UUID: resourcePack.UUID, Result: resourcePackResultDeclined},
			})
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigSelectKnownPacks:
			var knownPacks configclient.ConfigSelectKnownPacks
			if err = protocolpacket.Scan(p, &knownPacks); err != nil {
				return err
			}
			err = b.WritePacket(ctx, &configserver.ConfigSelectKnownPacks{})
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigCookieRequest:
			var request configclient.ConfigCookieRequest
			if err = protocolpacket.Scan(p, &request); err != nil {
				return err
			}
			err = b.WritePacket(ctx, &configserver.ConfigCookieResponse{
				CookieResponse: gameserver.CookieResponse{Key: request.Key},
			})
			if err != nil {
				return err
			}
		case packetid.ClientboundConfigCodeOfConduct:
			var codeOfConduct configclient.ConfigCodeOfConduct
			if err = protocolpacket.Scan(p, &codeOfConduct); err != nil {
				return err
			}
			// Minego auto-accepts the configuration code of conduct so the
			// connection can complete unattended instead of hanging in config.
			err = b.WritePacket(ctx, &configserver.ConfigAcceptCodeOfConduct{})
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
