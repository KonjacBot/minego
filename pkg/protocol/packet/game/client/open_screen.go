package client

import "github.com/KonjacBot/go-mc/chat"

//codec:gen
type OpenScreen struct {
	WindowID    int32 `mc:"VarInt"`
	WindowType  int32 `mc:"VarInt"`
	WindowTitle chat.Message
}
