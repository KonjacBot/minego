package client

import "github.com/KonjacBot/go-mc/chat"

//codec:gen
type Disconnect struct {
	Reason chat.Message
}
