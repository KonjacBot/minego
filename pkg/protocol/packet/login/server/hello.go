package server

import (
	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
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
