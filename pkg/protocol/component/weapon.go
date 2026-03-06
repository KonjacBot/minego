package component

//codec:gen
type Weapon struct {
	DamagePerAttack    int32   `mc:"VarInt"`
	DisableBlockingFor float32 // In seconds
}

func (*Weapon) ID() string {
	return "minecraft:weapon"
}
