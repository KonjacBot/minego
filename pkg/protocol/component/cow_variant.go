package component

//codec:gen
type CowVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*CowVariant) ID() string {
	return "minecraft:cow/variant"
}
