package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/level"
)

var _ ClientboundPacket = (*ForgetLevelChunk)(nil)

//codec:gen
type ForgetLevelChunk struct {
	Pos level.ChunkPos
}

func (ForgetLevelChunk) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundForgetLevelChunk
}
