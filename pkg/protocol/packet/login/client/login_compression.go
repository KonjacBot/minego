package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type LoginLoginCompression struct {
	Threshold int32 `mc:"VarInt"`
}

func (*LoginLoginCompression) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginCompression
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginCompression, func() ClientboundPacket {
		return &LoginLoginCompression{}
	})
}
