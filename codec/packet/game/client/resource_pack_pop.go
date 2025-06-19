package client

import (
	"github.com/google/uuid"
)

//codec:gen
type RemoveResourcePacket struct {
	HasUUID bool
	//opt:optional:HasUUID
	UUID uuid.UUID `mc:"UUID"`
}
