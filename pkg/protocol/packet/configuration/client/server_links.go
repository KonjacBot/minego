package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
)

type ConfigServerLinks struct {
	client.ServerLinks
}

func (*ConfigServerLinks) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigServerLinks
}

func init() {
	registerPacket(packetid.ClientboundConfigServerLinks, func() ClientboundPacket {
		return &ConfigServerLinks{}
	})
}
