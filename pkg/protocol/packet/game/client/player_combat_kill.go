package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type CombatDeath struct {
	PlayerID int32 `mc:"VarInt"`
	Message  chat.Message
}
