package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
)

type ConfigServerLinks struct {
	client.ServerLinks
}

func (ConfigServerLinks) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigServerLinks
}

func init() {
	registerPacket(packetid.ClientboundConfigServerLinks, func() ClientboundPacket {
		return &ConfigServerLinks{}
	})
}
