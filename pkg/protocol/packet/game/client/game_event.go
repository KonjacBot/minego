package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

var _ ClientboundPacket = (*GameEvent)(nil)

//codec:gen
type GameEvent struct {
	Event uint8
	Param float32
}

func (GameEvent) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundGameEvent
}
