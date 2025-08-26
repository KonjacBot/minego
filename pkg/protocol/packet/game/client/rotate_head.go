package client

import pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

//codec:gen
type SetHeadRotation struct {
	EntityID int32 `mc:"VarInt"`
	HeadYaw  pk.Angle
}
