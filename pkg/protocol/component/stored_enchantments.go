package component

//codec:gen
type StoredEnchantments struct {
	Enchantments []Enchantment
}

func (*StoredEnchantments) ID() string {
	return "minecraft:stored_enchantments"
}
