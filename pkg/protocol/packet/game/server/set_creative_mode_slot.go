package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
