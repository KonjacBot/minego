package server

import (
	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
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
