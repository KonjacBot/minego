package client

import (
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
)

//codec:gen
type ConfigDisconnect struct {
	Reason chat.Message
}

func (*ConfigDisconnect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigDisconnect
}

func init() {
	registerPacket(packetid.ClientboundConfigDisconnect, func() ClientboundPacket {
		return &ConfigDisconnect{}
	})
}
