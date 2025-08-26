package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PlayerCommand struct {
	EntityID  int32 `mc:"VarInt"`
	ActionID  int32 `mc:"VarInt"`
	JumpBoost int32 `mc:"VarInt"`
}

func (*PlayerCommand) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPlayerCommand
}

func init() {
	registerPacket(packetid.ServerboundPlayerCommand, func() ServerboundPacket {
		return &PlayerCommand{}
	})
}
