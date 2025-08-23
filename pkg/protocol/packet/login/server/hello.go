package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/google/uuid"
)

//codec:gen
type LoginHello struct {
	Name string
	UUID uuid.UUID `mc:"UUID"`
}

func (*LoginHello) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundLoginHello
}

func init() {
	registerPacket(packetid.ServerboundLoginHello, func() ServerboundPacket {
		return &LoginHello{}
	})
}
