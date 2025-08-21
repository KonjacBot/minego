package client

import "github.com/Tnze/go-mc/chat"

//codec:gen
type Disconnect struct {
	Reason chat.Message
}
