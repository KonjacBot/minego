package component

//codec:gen
type SheepColor struct {
	Color DyeColor
}

func (*SheepColor) ID() string {
	return "minecraft:sheep/color"
}
