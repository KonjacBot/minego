package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ConfigCookieRequest struct {
	Key string `mc:"Identifier"`
}

func (*ConfigCookieRequest) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigCookieRequest
}

func init() {
	registerPacket(packetid.ClientboundConfigCookieRequest, func() ClientboundPacket {
		return &ConfigCookieRequest{}
	})
}
