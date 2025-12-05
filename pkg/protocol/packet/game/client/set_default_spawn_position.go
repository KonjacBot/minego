package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type SetDefaultSpawnPosition struct {
	DimensionName pk.Identifier
	Location      pk.Position
	Yaw           float32
	Pitch         float32
}
