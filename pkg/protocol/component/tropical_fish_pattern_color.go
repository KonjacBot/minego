package component

//codec:gen
type TropicalFishPatternColor struct {
	Color DyeColor
}

func (*TropicalFishPatternColor) ID() string {
	return "minecraft:tropical_fish/pattern_color"
}
