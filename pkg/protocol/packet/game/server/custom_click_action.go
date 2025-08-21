package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type CustomClickAction struct {
	ID      string         `mc:"Identifier"`
	Payload nbt.RawMessage `mc:"NBT"`
}

func (CustomClickAction) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundCustomClickAction
}

func init() {
	registerPacket(packetid.ServerboundCustomClickAction, func() ServerboundPacket {
		return &CustomClickAction{}
	})
}
