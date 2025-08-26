package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
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
