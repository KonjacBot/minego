package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type SetTabListHeaderAndFooter struct {
	Header chat.Message
	Footer chat.Message
}
