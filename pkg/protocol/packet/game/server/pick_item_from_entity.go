package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PickItemFromEntity struct {
	EntityID    int32 `mc:"VarInt"`
	IncludeData bool
}

func (*PickItemFromEntity) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPickItemFromEntity
}

func init() {
	registerPacket(packetid.ServerboundPickItemFromEntity, func() ServerboundPacket {
		return &PickItemFromEntity{}
	})
}
