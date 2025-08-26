package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type ChangedSlot struct {
	Slot     int16
	SlotData slot.HashedSlot
}

//codec:gen
type ContainerClick struct {
	WindowID     int32 `mc:"VarInt"`
	StateID      int32 `mc:"VarInt"`
	Slot         int16
	Button       int8
	Mode         int32 `mc:"VarInt"`
	ChangedSlots []ChangedSlot
	CarriedSlot  slot.HashedSlot
}

func (*ContainerClick) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundContainerClick
}

func init() {
	registerPacket(packetid.ServerboundContainerClick, func() ServerboundPacket {
		return &ContainerClick{}
	})
}
