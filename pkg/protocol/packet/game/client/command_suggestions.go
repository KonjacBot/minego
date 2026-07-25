package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	chat "github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type CommandSuggestionsMatch struct {
	Match   string
	Tooltip pk.Option[chat.Message, *chat.Message]
}

//codec:gen
type CommandSuggestions struct {
	ID      int32 `mc:"VarInt"`
	Start   int32 `mc:"VarInt"`
	Length  int32 `mc:"VarInt"`
	Matches []CommandSuggestionsMatch
}
