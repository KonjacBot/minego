package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type CookieRequest struct {
	Key pk.Identifier
}
