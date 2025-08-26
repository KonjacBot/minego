package client

import "git.konjactw.dev/falloutBot/go-mc/chat"

//codec:gen
type CombatDeath struct {
	PlayerID int32 `mc:"VarInt"`
	Message  chat.Message
}
