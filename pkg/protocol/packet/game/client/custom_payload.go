package client

import "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type CustomPayload struct {
	Channel wire.Identifier
	Data    []byte `mc:"PluginMessageData"`
}
