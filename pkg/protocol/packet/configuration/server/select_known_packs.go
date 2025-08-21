package server

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/configuration/client"
	"github.com/Tnze/go-mc/data/packetid"
)

//codec:gen
type ConfigSelectKnownPacks struct {
	Packs []client.KnownPack
}

func (ConfigSelectKnownPacks) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigSelectKnownPacks
}

func init() {
	registerPacket(packetid.ServerboundConfigSelectKnownPacks, func() ServerboundPacket {
		return &ConfigSelectKnownPacks{}
	})
}
