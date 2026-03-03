package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type ConfigCustomClickAction struct {
	Action string         `mc:"Identifier"`
	Data   nbt.RawMessage `mc:"NBT"`
}

func (*ConfigCustomClickAction) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigCustomClickAction
}

func init() {
	registerPacket(packetid.ServerboundConfigCustomClickAction, func() ServerboundPacket {
		return &ConfigCustomClickAction{}
	})
}
