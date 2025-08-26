package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PlayerLoaded struct {
}

func (*PlayerLoaded) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPlayerLoaded
}

func init() {
	registerPacket(packetid.ServerboundPlayerLoaded, func() ServerboundPacket {
		return &PlayerLoaded{}
	})
}
