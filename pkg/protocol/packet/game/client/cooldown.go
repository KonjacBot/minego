package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	pk "github.com/Tnze/go-mc/net/packet"
)

var _ ClientboundPacket = (*Cooldown)(nil)
var _ pk.Field = (*Cooldown)(nil)

// CooldownPacket
//
//codec:gen
type Cooldown struct {
	CooldownGroup pk.Identifier `mc:"Identifier"`
	Duration      int32         `mc:"VarInt"`
}

func (Cooldown) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCooldown
}
