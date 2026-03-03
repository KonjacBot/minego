package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
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
