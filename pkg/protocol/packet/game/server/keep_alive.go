package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type KeepAlive struct {
	ID int64
}

func (KeepAlive) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundKeepAlive
}

func init() {
	registerPacket(packetid.ServerboundKeepAlive, func() ServerboundPacket {
		return &KeepAlive{}
	})
}
