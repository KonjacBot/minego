package client

//codec:gen
type HurtAnimation struct {
	EntityID int32 `mc:"VarInt"`
	Yaw      float32
}
