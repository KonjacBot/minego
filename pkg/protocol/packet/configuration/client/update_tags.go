package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
)

type ConfigUpdateTags struct {
	client.UpdateTags
}

func (*ConfigUpdateTags) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigUpdateTags
}

func init() {
	registerPacket(packetid.ClientboundConfigUpdateTags, func() ClientboundPacket {
		return &ConfigUpdateTags{}
	})
}
