package component

//codec:gen
type MapColor struct {
	Color int32 `mc:"Int"` // RGB components encoded as integer
}

func (*MapColor) ID() string {
	return "minecraft:map_color"
}
