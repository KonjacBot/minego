package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type LoginCustomQuery struct {
	MessageID int32  `mc:"VarInt"`
	Channel   string `mc:"Identifier"`
	Data      []byte `mc:"ByteArray"`
}

func (*LoginCustomQuery) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginCustomQuery
}

func init() {
	registerPacket(packetid.ClientboundLoginCustomQuery, func() ClientboundPacket {
		return &LoginCustomQuery{}
	})
}
