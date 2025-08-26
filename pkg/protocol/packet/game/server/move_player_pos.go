package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type MovePlayerPos struct {
	X, FeetY, Z float64
	Flags       int8
}

func (*MovePlayerPos) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundMovePlayerPos
}

func init() {
	registerPacket(packetid.ServerboundMovePlayerPos, func() ServerboundPacket {
		return &MovePlayerPos{}
	})
}
