package client

import (
	"git.konjactw.dev/patyhank/minego/codec/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
)

type ConfigResourcePackPop struct {
	client.RemoveResourcePack
}

func (ConfigResourcePackPop) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigResourcePackPop
}

func init() {
	registerPacket(packetid.ClientboundConfigResourcePackPop, func() ClientboundPacket {
		return &ConfigResourcePackPop{}
	})
}
