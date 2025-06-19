package client

import (
	"github.com/google/uuid"
)

//codec:gen
type RemoveResourcePack struct {
	HasUUID bool
	//opt:optional:HasUUID
	UUID uuid.UUID `mc:"UUID"`
}
