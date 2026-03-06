package component

//codec:gen
type VillagerVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*VillagerVariant) ID() string {
	return "minecraft:villager/variant"
}
