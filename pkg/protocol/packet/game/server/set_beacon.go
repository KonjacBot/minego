package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type SetBeacon struct {
	HasPrimaryEffect bool
	//opt:optional:HasPrimaryEffect
	PrimaryEffect      int32 `mc:"VarInt"`
	HasSecondaryEffect bool
	//opt:optional:HasSecondaryEffect
	SecondaryEffect int32 `mc:"VarInt"`
}

func (*SetBeacon) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetBeacon
}

func init() {
	registerPacket(packetid.ServerboundSetBeacon, func() ServerboundPacket {
		return &SetBeacon{}
	})
}
