package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

var _ ClientboundPacket = (*Cooldown)(nil)
var _ pk.Field = (*Cooldown)(nil)

// CooldownPacket
//
//codec:gen
type Cooldown struct {
	CooldownGroup wire.Identifier `mc:"Identifier"`
	Duration      int32           `mc:"VarInt"`
}

func (Cooldown) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundCooldown
}
