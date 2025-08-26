package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type PlayerAction struct {
	Status   int32 `mc:"VarInt"`
	Location pk.Position
	Face     int8
	Sequence int32 `mc:"VarInt"`
}

func (*PlayerAction) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPlayerAction
}

func init() {
	registerPacket(packetid.ServerboundPlayerAction, func() ServerboundPacket {
		return &PlayerAction{}
	})
}
