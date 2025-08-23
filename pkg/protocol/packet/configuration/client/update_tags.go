package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
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
