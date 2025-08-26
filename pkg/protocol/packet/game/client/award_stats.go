package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
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
