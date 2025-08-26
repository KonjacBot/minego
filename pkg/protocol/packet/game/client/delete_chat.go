package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*DeleteChat)(nil)
var _ packet.Field = (*DeleteChat)(nil)

// DeleteChatPacket
//
//codec:gen
type DeleteChat struct {
	MessageSignature []byte `mc:"ByteArray"`
}

func (DeleteChat) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDeleteChat
}
