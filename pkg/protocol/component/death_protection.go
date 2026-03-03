package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type DeathProtection struct {
	Effects []ConsumeEffect
}

func (*DeathProtection) Type() slot.ComponentID {
	return 32
}

func (*DeathProtection) ID() string {
	return "minecraft:death_protection"
}
