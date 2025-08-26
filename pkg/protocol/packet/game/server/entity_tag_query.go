package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type EntityTagQuery struct {
	TransactionID int32 `mc:"VarInt"`
	EntityID      int32 `mc:"VarInt"`
}

func (*EntityTagQuery) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundEntityTagQuery
}

func init() {
	registerPacket(packetid.ServerboundEntityTagQuery, func() ServerboundPacket {
		return &EntityTagQuery{}
	})
}
