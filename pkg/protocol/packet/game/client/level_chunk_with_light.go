package client

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/level"
)

var _ ClientboundPacket = (*LevelChunkWithLight)(nil)

type LevelChunkWithLight struct {
	Pos  level.ChunkPos
	Data *level.Chunk
}

func (c *LevelChunkWithLight) ReadFrom(r io.Reader) (n int64, err error) {
	temp, err := c.Pos.ReadFrom(r)
	if err != nil {
		return temp, err
	}
	c.Data = level.EmptyChunk(36)

	temp, err = (c.Data).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	return n, err
}

func (c LevelChunkWithLight) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = c.Pos.WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}

	temp, err = (*level.Chunk)(c.Data).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	return n, err
}
