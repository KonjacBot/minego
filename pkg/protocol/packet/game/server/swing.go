package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type Swing struct {
	Hand int32 `mc:"VarInt"`
}

func (*Swing) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSwing
}

func init() {
	registerPacket(packetid.ServerboundSwing, func() ServerboundPacket {
		return &Swing{}
	})
}
