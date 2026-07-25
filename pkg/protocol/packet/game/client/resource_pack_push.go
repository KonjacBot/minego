package client

import (
	"github.com/google/uuid"

	chat "github.com/KonjacBot/minego/pkg/protocol/wire"
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
