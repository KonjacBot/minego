package player

import "git.konjactw.dev/falloutBot/go-mc/chat"

type MessageEvent struct {
	Message chat.Message
}

func (m MessageEvent) EventID() string {
	return "player:message"
}
