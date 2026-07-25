package client

import "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type CookieRequest struct {
	Key wire.Identifier
}
