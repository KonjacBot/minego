package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

//codec:gen
type LoginLoginDisconnect struct {
	Reason chat.JsonMessage
}

func (*LoginLoginDisconnect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginDisconnect
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginDisconnect, func() ClientboundPacket {
		return &LoginLoginDisconnect{}
	})
}
