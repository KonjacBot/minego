package client

import (
	"github.com/google/uuid"
)

//codec:gen
type PlayerInfoRemove struct {
	UUIDs []uuid.UUID `mc:"UUID"`
}
