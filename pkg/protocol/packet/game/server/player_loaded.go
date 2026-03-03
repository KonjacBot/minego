package server

import "github.com/KonjacBot/go-mc/data/packetid"

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
