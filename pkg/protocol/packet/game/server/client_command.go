package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ClientCommand struct {
	Action int32 `mc:"VarInt"`
}

func (*ClientCommand) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundClientCommand
}

func init() {
	registerPacket(packetid.ServerboundClientCommand, func() ServerboundPacket {
		return &ClientCommand{}
	})
}
