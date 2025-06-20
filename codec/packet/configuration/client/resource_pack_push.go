package client

import (
	"git.konjactw.dev/patyhank/minego/codec/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
)

type ConfigResourcePackPush struct {
	client.AddResourcePack
}

func (ConfigResourcePackPush) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigResourcePackPush
}

func init() {
	registerPacket(packetid.ClientboundConfigResourcePackPush, func() ClientboundPacket {
		return &ConfigResourcePackPush{}
	})
}
