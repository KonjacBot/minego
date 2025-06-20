package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type ConfigFinishConfiguration struct {
}

func (ConfigFinishConfiguration) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigFinishConfiguration
}

func init() {
	registerPacket(packetid.ServerboundConfigFinishConfiguration, func() ServerboundPacket {
		return &ConfigFinishConfiguration{}
	})
}
