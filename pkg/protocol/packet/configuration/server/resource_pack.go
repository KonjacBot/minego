package server

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
	"github.com/Tnze/go-mc/data/packetid"
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
