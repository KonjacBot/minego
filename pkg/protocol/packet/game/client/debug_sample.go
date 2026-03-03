package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/net/packet"
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
