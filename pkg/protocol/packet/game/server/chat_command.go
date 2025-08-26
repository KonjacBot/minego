package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ChatCommand struct {
	Command string
}

func (*ChatCommand) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChatCommand
}

func init() {
	registerPacket(packetid.ServerboundChatCommand, func() ServerboundPacket {
		return &ChatCommand{}
	})
}
