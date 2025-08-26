package client

import "git.konjactw.dev/falloutBot/go-mc/chat"

//codec:gen
type SetTabListHeaderAndFooter struct {
	Header chat.Message
	Footer chat.Message
}
