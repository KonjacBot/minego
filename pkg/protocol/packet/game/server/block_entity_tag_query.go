package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type BlockEntityTagQuery struct {
	TransactionID int32 `mc:"VarInt"`
	Location      pk.Position
}

func (*BlockEntityTagQuery) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundBlockEntityTagQuery
}

func init() {
	registerPacket(packetid.ServerboundBlockEntityTagQuery, func() ServerboundPacket {
		return &BlockEntityTagQuery{}
	})
}
