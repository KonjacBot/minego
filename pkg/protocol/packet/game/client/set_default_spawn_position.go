package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type SetDefaultSpawnPosition struct {
	Location pk.Position
	Angle    float32
}
