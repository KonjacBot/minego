package component

//codec:gen
type ShulkerColor struct {
	Color DyeColor
}

func (*ShulkerColor) ID() string {
	return "minecraft:shulker/color"
}
