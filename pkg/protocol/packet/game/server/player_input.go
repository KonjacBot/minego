package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PlayerInput struct {
	Flags uint8
}

func (*PlayerInput) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPlayerInput
}

func init() {
	registerPacket(packetid.ServerboundPlayerInput, func() ServerboundPacket {
		return &PlayerInput{}
	})
}
