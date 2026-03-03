package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type EnchantmentGlintOverride struct {
	HasGlint bool
}

func (*EnchantmentGlintOverride) Type() slot.ComponentID {
	return 18
}

func (*EnchantmentGlintOverride) ID() string {
	return "minecraft:enchantment_glint_override"
}
