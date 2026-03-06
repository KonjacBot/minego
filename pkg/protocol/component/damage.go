package component

//codec:gen
type Damage struct {
	Damage int32 `mc:"VarInt"`
}

func (*Damage) ID() string {
	return "minecraft:damage"
}
