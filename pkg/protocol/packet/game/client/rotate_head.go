package client

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type SetHeadRotation struct {
	EntityID int32 `mc:"VarInt"`
	HeadYaw  pk.Angle
}
