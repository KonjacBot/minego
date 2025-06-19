package client

import (
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
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
