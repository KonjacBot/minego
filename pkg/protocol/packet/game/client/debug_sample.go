package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*DebugSample)(nil)
var _ packet.Field = (*DebugSample)(nil)

// DebugSamplePacket
//
//codec:gen
type DebugSample struct {
	Sample          []int64
	DebugSampleType int32 `mc:"VarInt"`
}

func (DebugSample) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDebugSample
}
