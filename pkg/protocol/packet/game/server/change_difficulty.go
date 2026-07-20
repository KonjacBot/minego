package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type ChangeDifficulty struct {
	Difficulty int32 `mc:"VarInt"`
}

func (*ChangeDifficulty) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChangeDifficulty
}

func init() {
	registerPacket(packetid.ServerboundChangeDifficulty, func() ServerboundPacket {
		return &ChangeDifficulty{}
	})
}
