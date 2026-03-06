package component

//codec:gen
type ZombieNautilusVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*ZombieNautilusVariant) ID() string {
	return "minecraft:zombie_nautilus/variant"
}
