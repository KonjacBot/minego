package component

//codec:gen
type FoxVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*FoxVariant) ID() string {
	return "minecraft:fox/variant"
}
