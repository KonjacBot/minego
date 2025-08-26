package client

import "git.konjactw.dev/falloutBot/go-mc/chat"

//codec:gen
type SystemChatMessage struct {
	Content chat.Message
	Overlay bool
}
