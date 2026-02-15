package client

import "git.konjactw.dev/patyhank/minego/pkg/protocol"

//codec:gen
type SetEntityVelocity struct {
	EntityID int32 `mc:"VarInt"`
	Velocity protocol.LpVec3
}
