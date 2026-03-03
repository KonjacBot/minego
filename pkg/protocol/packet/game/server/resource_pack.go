package server

import (
	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
)

//codec:gen
type ResourcePack struct {
	UUID   uuid.UUID `mc:"UUID"`
	Result int32     `mc:"VarInt"`
}

func (*ResourcePack) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundResourcePack
}

func init() {
	registerPacket(packetid.ServerboundResourcePack, func() ServerboundPacket {
		return &ResourcePack{}
	})
}
