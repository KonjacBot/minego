package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/yggdrasil/user"
	"github.com/google/uuid"
)

//codec:gen
type LoginLoginFinished struct {
	UUID       uuid.UUID `mc:"UUID"`
	Name       string
	Properties []user.Property
}

func (LoginLoginFinished) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginFinished
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginFinished, func() ClientboundPacket {
		return &LoginLoginFinished{}
	})
}
