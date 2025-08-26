package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*DamageEvent)(nil)
var _ pk.Field = (*DamageEvent)(nil)

//codec:gen
type DamageEventPos struct {
	X, Y, Z float64
}

// DamageEventPacket
//
//codec:gen
type DamageEvent struct {
	EntityID       int32 `mc:"VarInt"`
	SourceType     int32 `mc:"VarInt"`
	SourceCauseID  int32 `mc:"VarInt"`
	SourceDirectID int32 `mc:"VarInt"`
	SourcePos      pk.Option[DamageEventPos, *DamageEventPos]
}

func (DamageEvent) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDamageEvent
}
