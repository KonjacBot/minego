package client

import (
	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/chat"
)

//codec:gen
type AddResourcePack struct {
	UUID             uuid.UUID `mc:"UUID"`
	URL              string
	Hash             string
	Forced           bool
	HasPromptMessage bool
	//opt:optional:HasPromptMessage
	PromptMessage chat.Message
}
