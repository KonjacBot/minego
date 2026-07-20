package client

import "github.com/KonjacBot/go-mc/level"

type ChunkBiomeData struct {
	Pos  level.ChunkPos
	Data []byte `mc:"ByteArray"`
}

type ChunkBiomes struct {
	Chunks []ChunkBiomeData
}
