package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

//codec:gen
type Interact struct {
	EntityID int32 `mc:"VarInt"`
	Type     int32 `mc:"VarInt"`
	//opt:enum:Type:0
	InteractHand int32 `mc:"VarInt"`
	//opt:enum:Type:2
	InteractAtTargetX, InteractAtTargetY, InteractAtTargetZ float32
	//opt:enum:Type:2
	InteractAtHand  int32 `mc:"VarInt"`
	SneakKeyPressed bool
}

func (*Interact) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundInteract
}

func init() {
	registerPacket(packetid.ServerboundInteract, func() ServerboundPacket {
		return &Interact{}
	})
}
