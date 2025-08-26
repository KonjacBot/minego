package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

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
