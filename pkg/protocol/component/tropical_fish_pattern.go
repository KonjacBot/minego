package component

//codec:gen
type TropicalFishPattern struct {
	Pattern int32 `mc:"VarInt"`
}

func (*TropicalFishPattern) ID() string {
	return "minecraft:tropical_fish/pattern"
}
