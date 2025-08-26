package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type SetCarriedItem struct {
	Slot int16
}

func (*SetCarriedItem) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetCarriedItem
}

func init() {
	registerPacket(packetid.ServerboundSetCarriedItem, func() ServerboundPacket {
		return &SetCarriedItem{}
	})
}
