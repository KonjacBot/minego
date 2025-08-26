package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*Animate)(nil)
var _ packet.Field = (*Animate)(nil)

// AnimatePacket
// codec:gen
type Animate struct {
	EntityID int32 `mc:"VarInt"`
	Action   uint8
}

func (Animate) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAnimate
}
