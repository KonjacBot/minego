package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ContainerSlotStateChanged struct {
	SlotID   int32 `mc:"VarInt"`
	WindowID int32 `mc:"VarInt"`
	State    bool
}

func (*ContainerSlotStateChanged) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundContainerSlotStateChanged
}

func init() {
	registerPacket(packetid.ServerboundContainerSlotStateChanged, func() ServerboundPacket {
		return &ContainerSlotStateChanged{}
	})
}
