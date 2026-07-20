package component

type ContainerLoot struct {
	LootTable string `nbt:"loot_table"`
	Seed      int64  `nbt:"seed,omitempty"`
}

func (*ContainerLoot) ID() string {
	return "minecraft:container_loot"
}
