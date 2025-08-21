package client

import "github.com/Tnze/go-mc/chat"

//codec:gen
type OpenScreen struct {
	WindowID    int32 `mc:"VarInt"`
	WindowType  int32 `mc:"VarInt"`
	WindowTitle chat.Message
}
