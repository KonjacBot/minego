package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
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
