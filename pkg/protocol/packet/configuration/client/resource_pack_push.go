package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
)

type ConfigResourcePackPush struct {
	client.AddResourcePack
}

func (*ConfigResourcePackPush) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigResourcePackPush
}

func init() {
	registerPacket(packetid.ClientboundConfigResourcePackPush, func() ClientboundPacket {
		return &ConfigResourcePackPush{}
	})
}
