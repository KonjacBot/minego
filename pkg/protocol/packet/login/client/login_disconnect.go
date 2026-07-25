package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	chat "github.com/KonjacBot/minego/pkg/protocol/wire"
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
