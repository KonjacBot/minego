package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*Explode)(nil)

//codec:gen
type Vec3 struct {
	X, Y, Z float64
}

//codec:gen
type Explode struct {
	CenterX                 float64
	CenterY                 float64
	CenterZ                 float64
	PlayerKnockbackVelocity pk.Option[Vec3, *Vec3]
}

func (Explode) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundExplode
}
