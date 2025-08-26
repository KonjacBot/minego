package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type SignUpdate struct {
	Location                   pk.Position
	IsFrontText                bool
	Line1, Line2, Line3, Line4 string
}

func (*SignUpdate) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSignUpdate
}

func init() {
	registerPacket(packetid.ServerboundSignUpdate, func() ServerboundPacket {
		return &SignUpdate{}
	})
}
