package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

var _ ClientboundPacket = (*ContainerSetSlot)(nil)
var _ packet.Field = (*ContainerSetSlot)(nil)

// ContainerSetSlotPacket
//
//codec:gen
type ContainerSetSlot struct {
	ContainerID int32 `mc:"VarInt"`
	StateID     int32 `mc:"VarInt"`
	Slot        int16
	ItemStack   slot.Slot
}

func (ContainerSetSlot) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundContainerSetSlot
}
