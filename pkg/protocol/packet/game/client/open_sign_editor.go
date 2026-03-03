package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type OpenSignEditor struct {
	Location pk.Position
	Front    bool
}
