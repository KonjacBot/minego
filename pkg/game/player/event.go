package player

import "github.com/Tnze/go-mc/chat"

type MessageEvent struct {
	Message chat.Message
}

func (m MessageEvent) EventID() string {
	return "player:message"
}
