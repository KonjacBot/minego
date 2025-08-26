package server

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type TestInstanceBlockAction struct {
	Position pk.Position
	Action   int32 `mc:"VarInt"`
	IsTest   bool
	//opt:optional:IsTest
	Test                int32 `mc:"VarInt"`
	SizeX, SizeY, SizeZ int32 `mc:"VarInt"`
	Rotation            int32 `mc:"VarInt"`
	IgnoredEntities     bool
	Status              int32 `mc:"VarInt"`
	HasErrorMessage     bool
	//opt:optional:HasErrorMessage
	ErrorMessage chat.Message
}

func (*TestInstanceBlockAction) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundTestInstanceBlockAction
}

func init() {
	registerPacket(packetid.ServerboundTestInstanceBlockAction, func() ServerboundPacket {
		return &TestInstanceBlockAction{}
	})
}
