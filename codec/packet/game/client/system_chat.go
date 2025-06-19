package client

import "github.com/Tnze/go-mc/chat"

//codec:gen
type SystemChatMessage struct {
	Content chat.Message
	Overlay bool
}
