package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type SystemChatMessage struct {
	Content chat.Message
	Overlay bool
}
