package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

type ConfigResourcePack struct {
	server.ResourcePack
}

func (*ConfigResourcePack) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigResourcePack
}

func init() {
	registerPacket(packetid.ServerboundConfigResourcePack, func() ServerboundPacket {
		return &ConfigResourcePack{}
	})
}
