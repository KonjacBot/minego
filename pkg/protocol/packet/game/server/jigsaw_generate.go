package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
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
