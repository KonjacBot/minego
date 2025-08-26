package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type LoginHello struct {
	ServerID           string
	PublicKey          []byte `mc:"ByteArray"`
	VerifyToken        []byte `mc:"ByteArray"`
	ShouldAuthenticate bool
}

func (*LoginHello) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginHello
}

func init() {
	registerPacket(packetid.ClientboundLoginHello, func() ClientboundPacket {
		return &LoginHello{}
	})
}
