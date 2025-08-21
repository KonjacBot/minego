package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

var _ ClientboundPacket = (*CustomChatCompletions)(nil)
var _ packet.Field = (*CustomChatCompletions)(nil)

// CustomChatCompletionsPacket
//
//codec:gen
type CustomChatCompletions struct {
	Action  int32 `mc:"VarInt"`
	Entries []string
}

func (CustomChatCompletions) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCustomChatCompletions
}
