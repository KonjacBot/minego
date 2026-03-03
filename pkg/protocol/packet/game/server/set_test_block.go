package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type SetTestBlock struct {
	Position pk.Position
	Mode     int32 `mc:"VarInt"`
	Message  string
}

func (*SetTestBlock) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetTestBlock
}

func init() {
	registerPacket(packetid.ServerboundSetTestBlock, func() ServerboundPacket {
		return &SetTestBlock{}
	})
}
