package client

import "git.konjactw.dev/falloutBot/go-mc/chat"

//codec:gen
type Disconnect struct {
	Reason chat.Message
}
