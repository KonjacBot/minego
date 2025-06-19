package client

//codec:gen
type SetEntityVelocity struct {
	EntityID                        int32 `mc:"VarInt"`
	VelocityX, VelocityY, VelocityZ int16
}
