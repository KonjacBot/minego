package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type OpenScreen struct {
	WindowID    int32 `mc:"VarInt"`
	WindowType  int32 `mc:"VarInt"`
	WindowTitle chat.Message
}
