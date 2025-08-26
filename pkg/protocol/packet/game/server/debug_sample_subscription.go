package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type DebugSampleSubscription struct {
	SampleType int32 `mc:"VarInt"`
}

func (*DebugSampleSubscription) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundDebugSampleSubscription
}

func init() {
	registerPacket(packetid.ServerboundDebugSampleSubscription, func() ServerboundPacket {
		return &DebugSampleSubscription{}
	})
}
