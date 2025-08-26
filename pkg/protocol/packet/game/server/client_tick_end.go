package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ClientTickEnd struct {
}

func (*ClientTickEnd) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundClientTickEnd
}

func init() {
	registerPacket(packetid.ServerboundClientTickEnd, func() ServerboundPacket {
		return &ClientTickEnd{}
	})
}
