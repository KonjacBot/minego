package client

import (
	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*DisguisedChat)(nil)
var _ packet.Field = (*DisguisedChat)(nil)

// DisguisedChatPacket
//
//codec:gen
type DisguisedChat struct {
	Message  chat.Message
	ChatType []byte `mc:"ByteArray"`
}

func (DisguisedChat) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDisguisedChat
}
