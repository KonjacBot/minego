package component

//codec:gen
type Enchantable struct {
	Value int32 `mc:"VarInt"`
}

func (*Enchantable) ID() string {
	return "minecraft:enchantable"
}
