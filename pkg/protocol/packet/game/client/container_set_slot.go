package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
