package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

// codec:gen
type StatsData struct {
	CategoryID int32 `mc:"VarInt"`
	StatID     int32 `mc:"VarInt"`
	Value      int32 `mc:"VarInt"`
}

var _ ClientboundPacket = (*AwardStats)(nil)
var _ packet.Field = (*AwardStats)(nil)

// codec:gen
type AwardStats struct {
	Stats []StatsData
}

func (AwardStats) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAwardStats
}
