package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type FireworkExplosion struct {
	Explosion FireworkExplosionData
}

//codec:gen
type FireworkExplosionData struct {
	Shape      int32 `mc:"VarInt"`
	Colors     []int32
	FadeColors []int32
	HasTrail   bool
	HasTwinkle bool
}

func (*FireworkExplosion) Type() slot.ComponentID {
	return 59
}

func (*FireworkExplosion) ID() string {
	return "minecraft:firework_explosion"
}
