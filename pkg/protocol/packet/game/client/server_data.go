package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
)

//codec:gen
type ServerData struct {
	MOTD    chat.Message
	HasIcon bool
	//opt:optional:HasIcon
	Icon []int8 `mc:"Byte"`
}
