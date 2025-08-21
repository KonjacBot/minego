package server

import (
	"git.konjactw.dev/patyhank/minego/codec/packet/game/server"
	"github.com/Tnze/go-mc/data/packetid"
)

type ConfigCookieResponse struct {
	server.CookieResponse
}

func (ConfigCookieResponse) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigCookieResponse
}

func init() {
	registerPacket(packetid.ServerboundConfigCookieResponse, func() ServerboundPacket {
		return &ConfigCookieResponse{}
	})
}
