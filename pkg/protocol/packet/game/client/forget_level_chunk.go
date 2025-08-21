package client

import (
	"github.com/Tnze/go-mc/data/packetid"
)

var _ ClientboundPacket = (*ForgetLevelChunk)(nil)

//codec:gen
type ChunkPos struct {
	X, Z int32
}

//codec:gen
type ForgetLevelChunk struct {
	Pos ChunkPos
}

func (ForgetLevelChunk) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundForgetLevelChunk
}
