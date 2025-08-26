package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
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
