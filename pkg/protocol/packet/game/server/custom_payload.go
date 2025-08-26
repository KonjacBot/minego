package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type CustomPayload struct {
	Channel string `mc:"Identifier"`
	Data    []byte `mc:"ByteArray"`
}

func (*CustomPayload) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundCustomPayload
}

func init() {
	registerPacket(packetid.ServerboundCustomPayload, func() ServerboundPacket {
		return &CustomPayload{}
	})
}
