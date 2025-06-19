//codec:ignore
package server

import "github.com/Tnze/go-mc/data/packetid"

type ServerboundPacket interface {
	ServerboundPacketID() packetid.ServerboundPacketID
}

var ServerboundPackets = make(map[packetid.ServerboundPacketID]ServerboundPacket)
