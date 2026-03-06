package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type Pong struct {
	ID int32
}

func (*Pong) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPong
}

func init() {
	registerPacket(packetid.ServerboundPong, func() ServerboundPacket {
		return &Pong{}
	})
}
