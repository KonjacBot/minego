package client

//codec:gen
type SynchronizeVehiclePosition struct {
	EntityID                        int32 `mc:"VarInt"`
	X, Y, Z                         float64
	VelocityX, VelocityY, VelocityZ float64
	Yaw, Pitch                      float32
	Flags                           int32
	OnGround                        bool
}
