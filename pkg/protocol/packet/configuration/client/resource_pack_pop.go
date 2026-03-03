package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
)

type ConfigResourcePackPop struct {
	client.RemoveResourcePack
}

func (*ConfigResourcePackPop) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigResourcePackPop
}

func init() {
	registerPacket(packetid.ClientboundConfigResourcePackPop, func() ClientboundPacket {
		return &ConfigResourcePackPop{}
	})
}
