package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/configuration/client"
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
