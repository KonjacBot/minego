package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
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
