package component

//codec:gen
type MooshroomVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*MooshroomVariant) ID() string {
	return "minecraft:mooshroom/variant"
}
