package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type UseItemOn struct {
	Hand                      int32 `mc:"VarInt"`
	Location                  pk.Position
	Face                      int32 `mc:"VarInt"`
	CursorX, CursorY, CursorZ float32
	InsideBlock               bool
	WorldBorderHit            bool
	Sequence                  int32 `mc:"VarInt"`
}

func (*UseItemOn) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundUseItemOn
}

func init() {
	registerPacket(packetid.ServerboundUseItemOn, func() ServerboundPacket {
		return &UseItemOn{}
	})
}
