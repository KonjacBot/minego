package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type Fireworks struct {
	FlightDuration int32 `mc:"VarInt"`
	Explosions     []FireworkExplosionData
}

func (*Fireworks) Type() slot.ComponentID {
	return 60
}

func (*Fireworks) ID() string {
	return "minecraft:fireworks"
}
