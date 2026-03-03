package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type MovePlayerRot struct {
	XRot, YRot float32
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
