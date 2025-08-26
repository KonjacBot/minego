package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PingRequest struct {
	payload int64
}

func (*PingRequest) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPingRequest
}

func init() {
	registerPacket(packetid.ServerboundPingRequest, func() ServerboundPacket {
		return &PingRequest{}
	})
}
