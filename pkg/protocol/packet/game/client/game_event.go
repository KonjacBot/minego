package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
)

var _ ClientboundPacket = (*GameEvent)(nil)

const (
	// GameEventLevelChunksLoadStart tells the client that the initial chunk
	// stream is about to begin. Minecraft 26.2's level-load tracker does not
	// leave its WaitingForServer state until it receives this event.
	GameEventLevelChunksLoadStart uint8 = 13
)

//codec:gen
type GameEvent struct {
	Event uint8
	Param float32
}

func (GameEvent) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundGameEvent
}
