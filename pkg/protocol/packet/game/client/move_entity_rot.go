package client

//codec:gen
type UpdateEntityRotation struct {
	EntityID   int32 `mc:"VarInt"`
	Yaw, Pitch float32
	OnGround   bool
}
