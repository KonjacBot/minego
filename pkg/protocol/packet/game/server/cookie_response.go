package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type CookieResponse struct {
	Key        string `mc:"Identifier"`
	HasPayload bool
	//opt:optional:HasPayload
	Payload []byte `mc:"ByteArray"`
}

func (*CookieResponse) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundCookieResponse
}

func init() {
	registerPacket(packetid.ServerboundCookieResponse, func() ServerboundPacket {
		return &CookieResponse{}
	})
}
