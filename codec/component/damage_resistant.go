package component

import (
	"git.konjactw.dev/patyhank/minego/codec/slot"
	pk "github.com/Tnze/go-mc/net/packet"
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
