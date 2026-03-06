package component

//codec:gen
type CatCollar struct {
	Color DyeColor
}

func (*CatCollar) ID() string {
	return "minecraft:cat/collar"
}
