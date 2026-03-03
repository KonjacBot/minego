package client

import "github.com/KonjacBot/minego/pkg/protocol"

//codec:gen
type SetEntityVelocity struct {
	EntityID int32 `mc:"VarInt"`
	Velocity protocol.LpVec3
}
