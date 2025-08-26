package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type MoveVehicle struct {
	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool
}

func (*MoveVehicle) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundMoveVehicle
}

func init() {
	registerPacket(packetid.ServerboundMoveVehicle, func() ServerboundPacket {
		return &MoveVehicle{}
	})
}
