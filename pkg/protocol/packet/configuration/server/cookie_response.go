package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

type ConfigCookieResponse struct {
	server.CookieResponse
}

func (*ConfigCookieResponse) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigCookieResponse
}

func init() {
	registerPacket(packetid.ServerboundConfigCookieResponse, func() ServerboundPacket {
		return &ConfigCookieResponse{}
	})
}
