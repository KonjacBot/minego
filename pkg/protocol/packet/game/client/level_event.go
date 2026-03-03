package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*LevelEvent)(nil)

//codec:gen
type LevelEvent struct {
	Type        int32
	Pos         pk.Position
	Data        int32
	GlobalEvent bool
}

func (LevelEvent) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLevelEvent
}
