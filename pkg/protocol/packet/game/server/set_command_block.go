package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type SetCommandBlock struct {
	Location pk.Position
	Command  string
	Mode     int32 `mc:"VarInt"`
	Flags    int8
}

func (*SetCommandBlock) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetCommandBlock
}

func init() {
	registerPacket(packetid.ServerboundSetCommandBlock, func() ServerboundPacket {
		return &SetCommandBlock{}
	})
}
