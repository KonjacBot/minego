package client

import "github.com/Tnze/go-mc/chat"

//codec:gen
type CombatDeath struct {
	PlayerID int32 `mc:"VarInt"`
	Message  chat.Message
}
