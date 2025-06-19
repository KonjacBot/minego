package client

//codec:gen
type UpdateSectionsBlocks struct {
	ChunkSectionPosition int64
	Blocks               []int64 `mc:"VarLong"`
}

func (s *UpdateSectionsBlocks) SetSectionPos(x, y, z int32) {
	s.ChunkSectionPosition = ((int64(x) & 0x3FFFFF) << 42) | (int64(y) & 0xFFFFF) | ((int64(z) & 0x3FFFFF) << 20)
}

func (s UpdateSectionsBlocks) ToSectionPos() (x, y, z int32) {
	sectionX := int32(s.ChunkSectionPosition >> 42)
	sectionY := int32(s.ChunkSectionPosition << 44 >> 44)
	sectionZ := int32(s.ChunkSectionPosition >> 22 >> 42)
	return sectionX, sectionY, sectionZ
}

func (s *UpdateSectionsBlocks) AddBlock(x, y, z int, stateID int32) {
	s.Blocks = append(s.Blocks, int64(stateID)<<12|(int64(x)<<8|int64(y)<<4|int64(z)))
}

func (s UpdateSectionsBlocks) ParseBlocks() map[[3]int32]int32 {
	m := make(map[[3]int32]int32)
	for _, block := range s.Blocks {
		blockStateId := block >> 12
		blockLocalX := (block >> 8) & 0xF
		blockLocalY := block & 0xF
		blockLocalZ := (block >> 4) & 0xF

		m[[3]int32{int32(blockLocalX), int32(blockLocalY), int32(blockLocalZ)}] = int32(blockStateId)
	}
	return m
}
