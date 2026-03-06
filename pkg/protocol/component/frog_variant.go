package component

//codec:gen
type FrogVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*FrogVariant) ID() string {
	return "minecraft:frog/variant"
}
