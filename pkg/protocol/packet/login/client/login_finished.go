package client

import (
	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/yggdrasil/user"
)

//codec:gen
type LoginLoginFinished struct {
	UUID       uuid.UUID `mc:"UUID"`
	Name       string
	Properties []user.Property
}

func (*LoginLoginFinished) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginFinished
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginFinished, func() ClientboundPacket {
		return &LoginLoginFinished{}
	})
}
