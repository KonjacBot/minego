package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/google/uuid"
)

//codec:gen
type TeleportToEntity struct {
	TargetPlayer uuid.UUID `mc:"UUID"`
}

func (TeleportToEntity) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundTeleportToEntity
}

func init() {
	registerPacket(packetid.ServerboundTeleportToEntity, func() ServerboundPacket {
		return &TeleportToEntity{}
	})
}
