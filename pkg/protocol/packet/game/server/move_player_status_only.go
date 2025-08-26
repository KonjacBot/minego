package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type MovePlayerStatusOnly struct {
	Flags int8
}

func (*MovePlayerStatusOnly) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundMovePlayerStatusOnly
}

func init() {
	registerPacket(packetid.ServerboundMovePlayerStatusOnly, func() ServerboundPacket {
		return &MovePlayerStatusOnly{}
	})
}
