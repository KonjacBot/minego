package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
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
