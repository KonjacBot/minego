package client

import "github.com/KonjacBot/go-mc/chat"

//codec:gen
type SetTabListHeaderAndFooter struct {
	Header chat.Message
	Footer chat.Message
}
