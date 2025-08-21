//codec:ignore
package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	pk "github.com/Tnze/go-mc/net/packet"
)

type ServerboundPacket interface {
	pk.Field
	PacketID() packetid.ServerboundPacketID
}

type serverPacketCreator func() ServerboundPacket

var packetRegistry = make(map[packetid.ServerboundPacketID]serverPacketCreator)

func registerPacket(id packetid.ServerboundPacketID, creator serverPacketCreator) {
	packetRegistry[id] = creator
}
