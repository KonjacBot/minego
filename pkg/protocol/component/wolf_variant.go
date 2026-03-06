package component

//codec:gen
type WolfVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*WolfVariant) ID() string {
	return "minecraft:wolf/variant"
}
