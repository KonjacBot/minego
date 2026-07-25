package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type SetTitleText struct {
	Text chat.Message
}
