package client

import (
	"github.com/Tnze/go-mc/chat"
)

//codec:gen
type ServerData struct {
	MOTD    chat.Message
	HasIcon bool
	//opt:optional:HasIcon
	Icon []int8 `mc:"Byte"`
}
