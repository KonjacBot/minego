package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type DamageResistant struct {
	Types pk.Identifier // Tag specifying damage types
}

func (*DamageResistant) Type() slot.ComponentID {
	return 24
}

func (*DamageResistant) ID() string {
	return "minecraft:damage_resistant"
}
