package client

import pk "github.com/Tnze/go-mc/net/packet"

//codec:gen
type UpdateEntityPositionAndRotation struct {
	EntityID               int32 `mc:"VarInt"`
	DeltaX, DeltaY, DeltaZ int16
	Yaw, Pitch             pk.Angle
	OnGround               bool
}
