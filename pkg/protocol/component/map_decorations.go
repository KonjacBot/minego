package component

type MapDecorations struct {
	Decorations map[string]MapDecorationEntry
}

type MapDecorationEntry struct {
	Type     string  `nbt:"type"`
	X        float64 `nbt:"x"`
	Z        float64 `nbt:"z"`
	Rotation float32 `nbt:"rotation"`
}

func (*MapDecorations) ID() string {
	return "minecraft:map_decorations"
}
