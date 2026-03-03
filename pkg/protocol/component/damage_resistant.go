package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
