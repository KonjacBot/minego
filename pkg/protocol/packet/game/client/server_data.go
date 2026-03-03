package client

import (
	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type ServerData struct {
	MOTD    chat.Message
	HasIcon bool
	//opt:optional:HasIcon
	Icon pk.ByteArray
}
