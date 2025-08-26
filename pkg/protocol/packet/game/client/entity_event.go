package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*EntityEvent)(nil)
var _ packet.Field = (*EntityEvent)(nil)

// EntityEventPacket
//
//codec:gen
type EntityEvent struct {
	EntityID int32
	EventID  int8
}

func (EntityEvent) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundEntityEvent
}
