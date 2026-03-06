package component

//codec:gen
type PotDecorations struct {
	Decorations []int32 `mc:"PrefixedArray"`
}

func (*PotDecorations) ID() string {
	return "minecraft:pot_decorations"
}
