package component

//codec:gen
type MaxDamage struct {
	Damage int32 `mc:"VarInt"`
}

func (*MaxDamage) ID() string {
	return "minecraft:max_damage"
}
