package client

import "github.com/Tnze/go-mc/level"

var _ ClientboundPacket = (*LevelChunkWithLight)(nil)

//codec:gen
type LevelChunkWithLight struct {
	X    int32
	Z    int32
	Data level.Chunk
}
