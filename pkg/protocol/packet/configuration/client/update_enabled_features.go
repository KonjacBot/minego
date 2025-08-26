package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ConfigUpdateEnabledFeatures struct {
	Features []string `mc:"Identifier"`
}

func (*ConfigUpdateEnabledFeatures) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigUpdateEnabledFeatures
}

func init() {
	registerPacket(packetid.ClientboundConfigUpdateEnabledFeatures, func() ClientboundPacket {
		return &ConfigUpdateEnabledFeatures{}
	})
}
