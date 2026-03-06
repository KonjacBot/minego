package component

//codec:gen
type PigVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*PigVariant) ID() string {
	return "minecraft:pig/variant"
}
