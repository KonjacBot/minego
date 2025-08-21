package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type StoredEnchantments struct {
	Enchantments []Enchantment
}

func (*StoredEnchantments) Type() slot.ComponentID {
	return 34
}

func (*StoredEnchantments) ID() string {
	return "minecraft:stored_enchantments"
}
