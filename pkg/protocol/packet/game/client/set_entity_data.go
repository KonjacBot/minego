package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/metadata"
)

//codec:gen
type SetEntityMetadata struct {
	EntityID int32 `mc:"VarInt"`
	Metadata metadata.EntityMetadata
}
