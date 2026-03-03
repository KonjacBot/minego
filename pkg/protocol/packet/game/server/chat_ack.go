package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type ChatAck struct {
	MessageCount int32 `mc:"VarInt"`
}

func (*ChatAck) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChatAck
}

func init() {
	registerPacket(packetid.ServerboundChatAck, func() ServerboundPacket {
		return &ChatAck{}
	})
}
