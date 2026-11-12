package client

import (
	"github.com/KonjacBot/minego/v2/pkg/protocol/metadata"
)

//codec:gen
type SetEntityMetadata struct {
	EntityID int32 `mc:"VarInt"`
	Metadata metadata.EntityMetadata
}
