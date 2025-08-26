package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type Pong struct {
}

func (*Pong) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPong
}

func init() {
	registerPacket(packetid.ServerboundPong, func() ServerboundPacket {
		return &Pong{}
	})
}
