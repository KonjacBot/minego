package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type MinecartStep struct {
	X, Y, Z                         float64
	VelocityX, VelocityY, VelocityZ float64
	Yaw, Pitch                      pk.Angle
	Weight                          float32
}

//codec:gen
type MoveMinecartAlongTrack struct {
	EntityID int32 `mc:"VarInt"`
	Steps    []MinecartStep
}
