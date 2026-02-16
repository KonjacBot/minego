//codec:ignore
package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
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

func CreatePacket(id packetid.ServerboundPacketID) (ServerboundPacket, bool) {
	creator, ok := packetRegistry[id]
	if !ok {
		return nil, false
	}
	return creator(), true
}
