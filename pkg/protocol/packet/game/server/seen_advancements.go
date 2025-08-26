package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type SeenAdvancements struct {
	Action int32 `mc:"VarInt"`
	//opt:enum:Action:0
	TabID string `mc:"Identifier"`
}

func (*SeenAdvancements) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSeenAdvancements
}

func init() {
	registerPacket(packetid.ServerboundSeenAdvancements, func() ServerboundPacket {
		return &SeenAdvancements{}
	})
}
