package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type ClientboundPacket interface {
	pk.Field
	PacketID() packetid.ClientboundPacketID
}

type packetCreator func() ClientboundPacket

var packetRegistry = make(map[packetid.ClientboundPacketID]packetCreator)

func registerPacket(id packetid.ClientboundPacketID, creator packetCreator) {
	packetRegistry[id] = creator
}
