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
	ChatType BoundChatType
}

//codec:gen
type BoundChatType struct {
	// HolderID is the one-based wire ID of the chat_type registry entry.
	HolderID   int32 `mc:"VarInt"`
	Name       chat.Message
	TargetName packet.Option[chat.Message, *chat.Message]
}

func (DisguisedChat) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDisguisedChat
}
