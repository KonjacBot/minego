package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ChangeDifficulty struct {
	Difficulty uint8
}

func (*ChangeDifficulty) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChangeDifficulty
}

func init() {
	registerPacket(packetid.ServerboundChangeDifficulty, func() ServerboundPacket {
		return &ChangeDifficulty{}
	})
}
