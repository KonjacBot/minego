package client

import chat "github.com/KonjacBot/minego/pkg/protocol/wire"

//codec:gen
type SetSubtitleText struct {
	SubtitleText chat.Message
}
