package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type Disconnect struct {
	Reason chat.Message
}
