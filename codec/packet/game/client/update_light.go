package client

import "github.com/Tnze/go-mc/level"

//codec:gen
type UpdateLight struct {
	ChunkX, ChunkZ int32 `mc:"VarInt"`
	Data           level.LightData
}
