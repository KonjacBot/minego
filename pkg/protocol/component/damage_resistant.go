package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type DamageResistant struct {
	Types pk.Identifier // Tag specifying damage types
}

func (*DamageResistant) ID() string {
	return "minecraft:damage_resistant"
}
