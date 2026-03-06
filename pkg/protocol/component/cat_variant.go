package component

//codec:gen
type CatVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*CatVariant) ID() string {
	return "minecraft:cat/variant"
}
