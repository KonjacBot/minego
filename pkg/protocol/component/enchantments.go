package component

//codec:gen
type Enchantments struct {
	Enchantments []Enchantment
}

//codec:gen
type Enchantment struct {
	Type  int32 `mc:"VarInt"`
	Level int32 `mc:"VarInt"`
}

func (*Enchantments) ID() string {
	return "minecraft:enchantments"
}
