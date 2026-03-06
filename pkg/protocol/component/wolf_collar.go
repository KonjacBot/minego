package component

//codec:gen
type WolfCollar struct {
	Color DyeColor
}

func (*WolfCollar) ID() string {
	return "minecraft:wolf/collar"
}
