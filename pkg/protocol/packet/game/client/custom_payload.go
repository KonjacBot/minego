package client

import pk "github.com/Tnze/go-mc/net/packet"

//codec:gen
type CustomPayload struct {
	Channel pk.Identifier
	Data    []byte `mc:"ByteArray"`
}
