package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type TestInstanceBlockAction struct {
	Position pk.Position
	Action   int32 `mc:"VarInt"`
	Data     TestInstanceBlockData
}

func (*TestInstanceBlockAction) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundTestInstanceBlockAction
}

func init() {
	registerPacket(packetid.ServerboundTestInstanceBlockAction, func() ServerboundPacket {
		return &TestInstanceBlockAction{}
	})
}
