package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type SetActionBarText struct {
	Text chat.Message
}
