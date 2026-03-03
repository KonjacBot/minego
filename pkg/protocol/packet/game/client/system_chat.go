package client

import "github.com/KonjacBot/go-mc/chat"

//codec:gen
type SystemChatMessage struct {
	Content chat.Message
	Overlay bool
}
