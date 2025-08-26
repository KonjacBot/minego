package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
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
