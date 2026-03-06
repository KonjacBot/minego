package component

//codec:gen
type AxolotlVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*AxolotlVariant) ID() string {
	return "minecraft:axolotl/variant"
}
