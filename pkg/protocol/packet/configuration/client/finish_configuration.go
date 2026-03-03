package client

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type ConfigFinishConfiguration struct {
}

func (*ConfigFinishConfiguration) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigFinishConfiguration
}

func init() {
	registerPacket(packetid.ClientboundConfigFinishConfiguration, func() ClientboundPacket {
		return &ConfigFinishConfiguration{}
	})
}
