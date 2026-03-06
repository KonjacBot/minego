package component

//codec:gen
type HorseVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*HorseVariant) ID() string {
	return "minecraft:horse/variant"
}
