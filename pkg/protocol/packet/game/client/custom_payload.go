package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type CustomPayload struct {
	Channel pk.Identifier
	Data    []byte `mc:"ByteArray"`
}
