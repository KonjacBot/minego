package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type UpdateEntityRotation struct {
	EntityID   int32 `mc:"VarInt"`
	Yaw, Pitch pk.Angle
	OnGround   bool
}
