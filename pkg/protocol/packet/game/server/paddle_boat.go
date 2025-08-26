package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PaddleBoat struct {
	LeftTurning, RightTurning bool
}

func (*PaddleBoat) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPaddleBoat
}

func init() {
	registerPacket(packetid.ServerboundPaddleBoat, func() ServerboundPacket {
		return &PaddleBoat{}
	})
}
