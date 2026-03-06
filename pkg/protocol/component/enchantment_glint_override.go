package component

//codec:gen
type EnchantmentGlintOverride struct {
	HasGlint bool
}

func (*EnchantmentGlintOverride) ID() string {
	return "minecraft:enchantment_glint_override"
}
