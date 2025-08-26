package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type CookieRequest struct {
	Key pk.Identifier
}
