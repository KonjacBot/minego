package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type UseItem struct {
	Hand       int32 `mc:"VarInt"`
	Sequence   int32 `mc:"VarInt"`
	Yaw, Pitch float32
}

func (*UseItem) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundUseItem
}

func init() {
	registerPacket(packetid.ServerboundUseItem, func() ServerboundPacket {
		return &UseItem{}
	})
}
