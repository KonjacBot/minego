package server

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type PingRequest struct {
	Payload int64
}

func (*PingRequest) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPingRequest
}

func init() {
	registerPacket(packetid.ServerboundPingRequest, func() ServerboundPacket {
		return &PingRequest{}
	})
}
