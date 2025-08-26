package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ConfigPong struct {
	ID int32
}

func (*ConfigPong) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigPong
}

func init() {
	registerPacket(packetid.ServerboundConfigPong, func() ServerboundPacket {
		return &ConfigPong{}
	})
}
