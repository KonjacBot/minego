package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/configuration/client"
)

//codec:gen
type ConfigSelectKnownPacks struct {
	Packs []client.KnownPack
}

func (*ConfigSelectKnownPacks) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigSelectKnownPacks
}

func init() {
	registerPacket(packetid.ServerboundConfigSelectKnownPacks, func() ServerboundPacket {
		return &ConfigSelectKnownPacks{}
	})
}
