package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PlayerAbilities struct {
	Flags int8
}

func (*PlayerAbilities) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPlayerAbilities
}

func init() {
	registerPacket(packetid.ServerboundPlayerAbilities, func() ServerboundPacket {
		return &PlayerAbilities{}
	})
}
