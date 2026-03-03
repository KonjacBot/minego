package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type SetDefaultSpawnPosition struct {
	DimensionName pk.Identifier
	Location      pk.Position
	Yaw           float32
	Pitch         float32
}
