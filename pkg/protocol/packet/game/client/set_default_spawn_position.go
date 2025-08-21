package client

import pk "github.com/Tnze/go-mc/net/packet"

//codec:gen
type SetDefaultSpawnPosition struct {
	Location pk.Position
	Angle    float32
}
