package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type SetJigsawBlock struct {
	Location          pk.Position
	Name              string `mc:"Identifier"`
	Target            string `mc:"Identifier"`
	Pool              string `mc:"Identifier"`
	FinalState        string
	JointType         string
	SelectionPriority int32 `mc:"VarInt"`
	PlacementPriority int32 `mc:"VarInt"`
}

func (*SetJigsawBlock) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetJigsawBlock
}

func init() {
	registerPacket(packetid.ServerboundSetJigsawBlock, func() ServerboundPacket {
		return &SetJigsawBlock{}
	})
}
