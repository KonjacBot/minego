package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type SetCreativeModeSlot struct {
	Slot        int16
	ClickedItem slot.Slot
}

func (*SetCreativeModeSlot) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetCreativeModeSlot
}

func init() {
	registerPacket(packetid.ServerboundSetCreativeModeSlot, func() ServerboundPacket {
		return &SetCreativeModeSlot{}
	})
}
