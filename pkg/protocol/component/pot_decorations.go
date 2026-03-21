package component

//codec:gen
type PotDecorations struct {
	Decorations []int32 `mc:"VarInt"`
}

func (*PotDecorations) ID() string {
	return "minecraft:pot_decorations"
}
