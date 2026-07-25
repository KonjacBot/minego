package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type DamageResistant struct {
	Types wire.IDSet // HolderSet of damage types
}

func (*DamageResistant) ID() string {
	return "minecraft:damage_resistant"
}
