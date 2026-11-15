package keeper

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
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

// clientInformation is what the keeper tells the server about its client, the
// same values minego's own client sends, so a kept connection looks the same
// as a direct one.
var clientInformation = gameserver.ClientInformation{
	Location:           "zh_TW",
	ViewDistance:       16,
	ChatColor:          true,
	DisplayedSkinParts: 127,
	AllowListing:       true,
}

// dial opens the server connection and logs in. The access token is used for
// the session join only and dropped afterwards.
func (s *session) dial(ctx context.Context) error {
	srv := s.srv
	host, port, err := splitAddress(srv.Address)
	if err != nil {
		return err
	}
	var dialer mcnet.MCDialer = &mcnet.DefaultDialer
	if srv.Proxy != nil {
		if dialer, err = srv.Proxy.Dialer(); err != nil {
			return err
		}
	}
	conn, err := dialer.DialMCContext(ctx, srv.Address)
	if err != nil {
		return fmt.Errorf("連線到伺服器: %w", err)
	}
	ok := false
	defer func() {
		if !ok {
			_ = conn.Close()
		}
	}()
	stop := context.AfterFunc(ctx, func() { _ = conn.Socket.SetDeadline(time.Now()) })
	defer stop()

	if srv.FakeHost != "" {
		host = srv.FakeHost
	}
	if err := conn.WritePacket(pk.Marshal(0,
		pk.VarInt(protocol.ProtocolVersion), pk.String(host), pk.UnsignedShort(port), pk.VarInt(2),
	)); err != nil {
		return fmt.Errorf("握手: %w", err)
	}
	var provider auth.Provider = &auth.OfflineAuth{Username: srv.Name}
	if srv.AccessToken != "" {
		provider = &auth.OnlineAuth{AccessToken: srv.AccessToken, Profile: auth.Profile{Name: srv.Name, UUID: srv.UUID}}
	}
	login := &auth.Auth{Conn: conn, Provider: provider, ReadIdleTimeout: s.k.readIdle}
	if err := login.HandleLogin(ctx); err != nil {
		return fmt.Errorf("登入: %w", err)
	}
	s.srv.AccessToken = ""
	s.finished = login.Finished
	s.up = conn
	if err := s.configure(ctx); err != nil {
		return err
	}
	ok = true
	return nil
}

func splitAddress(addr string) (string, uint16, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		var addrErr *net.AddrError
		if errors.As(err, &addrErr) && addrErr.Err == "missing port in address" {
			return addr, 25565, nil
		}
		return "", 0, fmt.Errorf("伺服器位址 %q: %w", addr, err)
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return "", 0, fmt.Errorf("伺服器位址 %q: %w", addr, err)
	}
	return host, uint16(port), nil
}

// configure runs the configuration phase against the server, answering what a
// client must answer and recording the rest for the core's own configuration.
// The keeper always configures by itself, also when a core is attached: the
// core is then taken through a replay of it (see core.go).
func (s *session) configure(ctx context.Context) error {
	var recorded []pk.Packet
	information := configserver.ConfigClientInformation{ClientInformation: clientInformation}
	if err := s.writeUp(pk.Marshal(packetid.ServerboundConfigClientInformation, &information)); err != nil {
		return err
	}
	for {
		p, err := s.readUp(ctx)
		if err != nil {
			return fmt.Errorf("設定階段: %w", err)
		}
		switch packetid.ClientboundPacketID(p.ID) {
		case packetid.ClientboundConfigDisconnect:
			var reason chat.Message
			_ = p.Scan(&reason)
			return &KickError{Reason: reason.ClearString()}
		case packetid.ClientboundConfigFinishConfiguration:
			if err := s.writeUp(pk.Marshal(packetid.ServerboundConfigFinishConfiguration)); err != nil {
				return err
			}
			s.mu.Lock()
			s.st.config = recorded
			s.mu.Unlock()
			return nil
		case packetid.ClientboundConfigKeepAlive:
			var id pk.Long
			if err := p.Scan(&id); err != nil {
				return err
			}
			if err := s.writeUp(pk.Marshal(packetid.ServerboundConfigKeepAlive, id)); err != nil {
				return err
			}
		case packetid.ClientboundConfigPing:
			var id pk.Int
			if err := p.Scan(&id); err != nil {
				return err
			}
			if err := s.writeUp(pk.Marshal(packetid.ServerboundConfigPong, id)); err != nil {
				return err
			}
		case packetid.ClientboundConfigSelectKnownPacks:
			var offered configclient.ConfigSelectKnownPacks
			if err := p.Scan(&offered); err != nil {
				return err
			}
			response := configserver.ConfigSelectKnownPacks{Packs: configclient.SupportedKnownPacks(offered.KnownPacks)}
			if err := s.writeUp(pk.Marshal(packetid.ServerboundConfigSelectKnownPacks, &response)); err != nil {
				return err
			}
			// The core answers the replay the same way (same protocol version), so the
			// registry data recorded after it matches what it expects.
			recorded = append(recorded, p)
		case packetid.ClientboundConfigStoreCookie:
			var stored configclient.ConfigStoreCookie
			if err := p.Scan(&stored); err != nil {
				return err
			}
			s.storeCookie(stored.Key, stored.Payload)
		case packetid.ClientboundConfigCookieRequest:
			var request configclient.ConfigCookieRequest
			if err := p.Scan(&request); err != nil {
				return err
			}
			response := configserver.ConfigCookieResponse{CookieResponse: s.cookieResponse(request.Key)}
			if err := s.writeUp(pk.Marshal(response.PacketID(), &response)); err != nil {
				return err
			}
		case packetid.ClientboundConfigResourcePackPush:
			var offered configclient.ConfigResourcePackPush
			if err := p.Scan(&offered); err != nil {
				return err
			}
			for _, result := range []int32{3, 0} {
				response := configserver.ConfigResourcePack{ResourcePack: gameserver.ResourcePack{UUID: offered.UUID, Result: result}}
				if err := s.writeUp(pk.Marshal(response.PacketID(), &response)); err != nil {
					return err
				}
			}
		default:
			recorded = append(recorded, p)
		}
	}
}

// KickError is the server ending the connection with a reason.
type KickError struct{ Reason string }

func (e *KickError) Error() string { return "被伺服器中斷連線: " + e.Reason }

// readUp reads one packet from the server, giving up after the idle timeout.
func (s *session) readUp(ctx context.Context) (pk.Packet, error) {
	if err := ctx.Err(); err != nil {
		return pk.Packet{}, err
	}
	_ = s.up.Socket.SetReadDeadline(time.Now().Add(s.k.readIdle))
	var p pk.Packet
	if err := s.up.ReadPacket(&p); err != nil {
		if ctx.Err() != nil {
			return pk.Packet{}, ctx.Err()
		}
		return pk.Packet{}, err
	}
	return p, nil
}

// writeUp sends one packet to the server. Writes come from the read loop
// (answers) and from the core's forwarder, so they are serialized here.
func (s *session) writeUp(p pk.Packet) error {
	s.upMu.Lock()
	defer s.upMu.Unlock()
	_ = s.up.Socket.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := s.up.WritePacket(p); err != nil {
		return fmt.Errorf("送封包到伺服器: %w", err)
	}
	return nil
}

func (s *session) storeCookie(key string, payload []int8) {
	s.cookieMu.Lock()
	defer s.cookieMu.Unlock()
	s.cookies[key] = append([]int8(nil), payload...)
}

func (s *session) cookieResponse(key string) gameserver.CookieResponse {
	s.cookieMu.Lock()
	defer s.cookieMu.Unlock()
	payload, ok := s.cookies[key]
	return gameserver.CookieResponse{Key: key, HasPayload: ok, Payload: append([]int8(nil), payload...)}
}
