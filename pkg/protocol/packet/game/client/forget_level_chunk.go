package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level"
)

var _ ClientboundPacket = (*ForgetLevelChunk)(nil)

//codec:gen
type ForgetLevelChunk struct {
	Pos level.ChunkPos
}

func (ForgetLevelChunk) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundForgetLevelChunk
}
