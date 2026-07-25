package client

import (
	"errors"
	"fmt"
	"io"

	"github.com/KonjacBot/go-mc/level"
	"github.com/KonjacBot/go-mc/level/block"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

var _ ClientboundPacket = (*LevelChunkWithLight)(nil)

type LevelChunkWithLight struct {
	Pos  level.ChunkPos
	Data *level.Chunk
}

func (c *LevelChunkWithLight) ReadFrom(r io.Reader) (n int64, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("decode level chunk panicked: %v", recovered)
		}
	}()

	temp, err := c.Pos.ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	c.Data = level.EmptyChunk(36)

	temp, err = readChunkData(r, c.Data)
	n += temp
	if err != nil {
		return n, err
	}
	return n, nil
}

func readChunkData(r io.Reader, chunk *level.Chunk) (n int64, err error) {
	var heightMapCount pk.VarInt
	temp, err := heightMapCount.ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if heightMapCount < 0 || heightMapCount > 64 {
		return n, fmt.Errorf("invalid height map count %d", heightMapCount)
	}
	if heightMapCount > 0 {
		chunk.HeightMaps = make(level.HeightMaps, int(heightMapCount))
	} else {
		chunk.HeightMaps = nil
	}
	for i := range chunk.HeightMaps {
		temp, err = (*pk.VarInt)(&chunk.HeightMaps[i].Type).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = wire.Array(&chunk.HeightMaps[i].Data).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}

	var sectionData wire.ByteArray
	temp, err = sectionData.ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if err = chunk.PutData(sectionData); err != nil {
		return n, err
	}
	if err = validateChunkStates(chunk); err != nil {
		return n, err
	}

	var blockEntityCount pk.VarInt
	temp, err = blockEntityCount.ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if blockEntityCount < 0 || blockEntityCount > wire.MaxCollectionEntries {
		return n, fmt.Errorf("invalid block entity count %d", blockEntityCount)
	}
	if blockEntityCount > 0 {
		chunk.BlockEntity = make([]level.BlockEntity, int(blockEntityCount))
	} else {
		chunk.BlockEntity = nil
	}
	for i := range chunk.BlockEntity {
		entity := &chunk.BlockEntity[i]
		temp, err = (*pk.Byte)(&entity.XZ).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = (*pk.Short)(&entity.Y).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = (*pk.VarInt)(&entity.Type).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
		temp, err = wire.NBT(&entity.Data).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}

	maxLightWords := (len(chunk.Sections) + 2 + 63) / 64
	for range 4 {
		temp, err = readLightMask(r, maxLightWords)
		n += temp
		if err != nil {
			return n, err
		}
	}
	for range 2 {
		temp, err = readLightArrays(r, len(chunk.Sections)+2)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func readLightMask(r io.Reader, maxWords int) (n int64, err error) {
	var length pk.VarInt
	n, err = length.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if length < 0 || int(length) > maxWords {
		return n, fmt.Errorf("invalid light mask length %d", length)
	}
	for range int(length) {
		var word pk.Long
		read, readErr := word.ReadFrom(r)
		n += read
		if readErr != nil {
			return n, readErr
		}
	}
	return n, nil
}

func readLightArrays(r io.Reader, maxArrays int) (n int64, err error) {
	var count pk.VarInt
	n, err = count.ReadFrom(r)
	if err != nil {
		return n, err
	}
	if count < 0 || int(count) > maxArrays {
		return n, fmt.Errorf("invalid light array count %d", count)
	}
	for range int(count) {
		var data wire.ByteArray
		read, readErr := data.ReadFrom(r)
		n += read
		if readErr != nil {
			return n, readErr
		}
		if len(data) != 2048 {
			return n, fmt.Errorf("invalid light array length %d", len(data))
		}
	}
	return n, nil
}

func validateChunkStates(chunk *level.Chunk) error {
	for sectionIndex := range chunk.Sections {
		section := &chunk.Sections[sectionIndex]
		for blockIndex := 0; blockIndex < 16*16*16; blockIndex++ {
			state := section.GetBlock(blockIndex)
			if state < 0 || int(state) >= len(block.StateList) {
				return fmt.Errorf("section %d contains invalid block state %d", sectionIndex, state)
			}
		}
	}
	return nil
}

func (c LevelChunkWithLight) WriteTo(w io.Writer) (n int64, err error) {
	if c.Data == nil {
		return 0, errors.New("level chunk data is nil")
	}
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
