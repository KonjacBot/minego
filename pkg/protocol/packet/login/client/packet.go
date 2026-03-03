package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
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
