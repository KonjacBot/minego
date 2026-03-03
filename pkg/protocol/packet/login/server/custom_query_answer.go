package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type LoginCustomQueryAnswer struct {
	MessageID int32 `mc:"VarInt"`
	HasData   bool
	//opt:optional:HasData
	Data []byte `mc:"PluginMessageData"`
}

func (*LoginCustomQueryAnswer) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundLoginCustomQueryAnswer
}

func init() {
	registerPacket(packetid.ServerboundLoginCustomQueryAnswer, func() ServerboundPacket {
		return &LoginCustomQueryAnswer{}
	})
}
