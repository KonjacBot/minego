package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
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
