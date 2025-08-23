package server

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/data/packetid"
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
