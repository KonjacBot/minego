package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type MovePlayerPosRot struct {
	X, FeetY, Z float64
	Yaw, Pitch  float32
	Flags       int8
}

func (*MovePlayerPosRot) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundMovePlayerPosRot
}

func init() {
	registerPacket(packetid.ServerboundMovePlayerPosRot, func() ServerboundPacket {
		return &MovePlayerPosRot{}
	})
}
