package component

//codec:gen
type BaseColor struct {
	Color DyeColor
}

func (*BaseColor) ID() string {
	return "minecraft:base_color"
}
