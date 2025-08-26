package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type CustomPayload struct {
	Channel pk.Identifier
	Data    []byte `mc:"ByteArray"`
}
