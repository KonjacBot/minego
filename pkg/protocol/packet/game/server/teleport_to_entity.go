package server

import (
	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
)

//codec:gen
type TeleportToEntity struct {
	TargetPlayer uuid.UUID `mc:"UUID"`
}

func (*TeleportToEntity) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundTeleportToEntity
}

func init() {
	registerPacket(packetid.ServerboundTeleportToEntity, func() ServerboundPacket {
		return &TeleportToEntity{}
	})
}
