package client

import (
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
)

//codec:gen
type LoginLoginDisconnect struct {
	Reason chat.JsonMessage
}

func (LoginLoginDisconnect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginDisconnect
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginDisconnect, func() ClientboundPacket {
		return &LoginLoginDisconnect{}
	})
}
