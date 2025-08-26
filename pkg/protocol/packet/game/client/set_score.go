package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/nbt"
)

//codec:gen
type ScoreNumberFormat struct {
	NumberFormat int32
	//opt:enum:NumberFormat:1
	StyledTag nbt.RawMessage `mc:"NBT"`
	//opt:enum:NumberFormat:2
	Content chat.Message
}

//codec:gen
type UpdateScore struct {
	EntityName     string
	ObjectiveName  string
	Value          int32 `mc:"VarInt"`
	HasDisplayName bool
	//opt:optional:HasDisplayName
	DisplayName    chat.Message
	HasScoreFormat bool
	//opt:optional:HasScoreFormat
	NumberFormat ScoreNumberFormat
}
