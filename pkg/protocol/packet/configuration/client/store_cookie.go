package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ConfigStoreCookie struct {
	Key     string `mc:"Identifier"`
	Payload []int8
}

func (*ConfigStoreCookie) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigStoreCookie
}

func init() {
	registerPacket(packetid.ClientboundConfigStoreCookie, func() ClientboundPacket {
		return &ConfigStoreCookie{}
	})
}
