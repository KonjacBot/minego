package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type DebugSampleSubscription struct {
	SampleType int32 `mc:"VarInt"`
}

func (*DebugSampleSubscription) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundDebugSubscriptionRequest
}

func init() {
	registerPacket(packetid.ServerboundDebugSubscriptionRequest, func() ServerboundPacket {
		return &DebugSampleSubscription{}
	})
}
