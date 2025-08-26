package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

// codec:gen
type BlockEvent struct {
	Position  packet.Position
	EventType uint8
	Data      uint8
	Block     int32 `mc:"VarInt"`
}

func (BlockEvent) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockEvent
}
