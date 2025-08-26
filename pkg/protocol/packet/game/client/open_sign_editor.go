package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type OpenSignEditor struct {
	Location pk.Position
	Front    bool
}
