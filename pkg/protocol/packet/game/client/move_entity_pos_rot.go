package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type UpdateEntityPositionAndRotation struct {
	EntityID               int32 `mc:"VarInt"`
	DeltaX, DeltaY, DeltaZ int16
	Yaw, Pitch             pk.Angle
	OnGround               bool
}
