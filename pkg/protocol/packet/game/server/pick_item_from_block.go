package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type PickItemFromBlock struct {
	Location    pk.Position
	IncludeData bool
}

func (*PickItemFromBlock) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPickItemFromBlock
}

func init() {
	registerPacket(packetid.ServerboundPickItemFromBlock, func() ServerboundPacket {
		return &PickItemFromBlock{}
	})
}
