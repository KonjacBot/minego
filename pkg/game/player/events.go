package player

import "github.com/KonjacBot/go-mc/chat"

type MessageEvent struct {
	Message chat.Message
}

func (m MessageEvent) EventID() string {
	return "player:message"
}
