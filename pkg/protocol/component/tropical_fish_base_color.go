package component

//codec:gen
type TropicalFishBaseColor struct {
	Color DyeColor
}

func (*TropicalFishBaseColor) ID() string {
	return "minecraft:tropical_fish/base_color"
}
