package client

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type ServerData struct {
	MOTD    wire.Message
	HasIcon bool
	//opt:optional:HasIcon
	Icon wire.ByteArray
}
