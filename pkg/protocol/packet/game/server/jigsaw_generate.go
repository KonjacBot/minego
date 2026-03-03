package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type JigsawGenerate struct {
	Location    pk.Position
	Levels      int32 `mc:"VarInt"`
	KeepJigsaws bool
}

func (*JigsawGenerate) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundJigsawGenerate
}

func init() {
	registerPacket(packetid.ServerboundJigsawGenerate, func() ServerboundPacket {
		return &JigsawGenerate{}
	})
}
