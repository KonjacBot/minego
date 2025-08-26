package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type MovePlayerRot struct {
	Yaw, Pitch float32
	Flags      int8
}

func (*MovePlayerRot) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundMovePlayerRot
}

func init() {
	registerPacket(packetid.ServerboundMovePlayerRot, func() ServerboundPacket {
		return &MovePlayerRot{}
	})
}
