package client

import "git.konjactw.dev/falloutBot/go-mc/level"

//codec:gen
type UpdateLight struct {
	ChunkX, ChunkZ int32 `mc:"VarInt"`
	Data           level.LightData
}
