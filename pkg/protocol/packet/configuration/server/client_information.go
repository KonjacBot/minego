package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
)

type ConfigClientInformation struct {
	server.ClientInformation
}

func (*ConfigClientInformation) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigClientInformation
}

func init() {
	registerPacket(packetid.ServerboundConfigClientInformation, func() ServerboundPacket {
		return &ConfigClientInformation{}
	})
}
