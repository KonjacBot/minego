package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type LoginLoginAcknowledged struct {
}

func (*LoginLoginAcknowledged) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundLoginLoginAcknowledged
}

func init() {
	registerPacket(packetid.ServerboundLoginLoginAcknowledged, func() ServerboundPacket {
		return &LoginLoginAcknowledged{}
	})
}
