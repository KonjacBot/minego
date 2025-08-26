package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type CookieResponse struct {
	Key        string `mc:"Identifier"`
	HasPayload bool
	//opt:optional:HasPayload
	Payload []int8 `mc:"Byte"`
}

func (*CookieResponse) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundCookieResponse
}

func init() {
	registerPacket(packetid.ServerboundCookieResponse, func() ServerboundPacket {
		return &CookieResponse{}
	})
}
