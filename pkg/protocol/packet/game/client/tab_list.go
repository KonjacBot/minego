package client

import "github.com/Tnze/go-mc/chat"

//codec:gen
type SetTabListHeaderAndFooter struct {
	Header chat.Message
	Footer chat.Message
}
