package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type ConfigurationAcknowledged struct {
}

func (*ConfigurationAcknowledged) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigurationAcknowledged
}

func init() {
	registerPacket(packetid.ServerboundConfigurationAcknowledged, func() ServerboundPacket {
		return &ConfigurationAcknowledged{}
	})
}
