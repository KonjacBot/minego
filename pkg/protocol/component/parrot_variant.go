package component

//codec:gen
type ParrotVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*ParrotVariant) ID() string {
	return "minecraft:parrot/variant"
}
