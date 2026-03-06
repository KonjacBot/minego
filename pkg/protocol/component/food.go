package component

//codec:gen
type Food struct {
	Nutrition          int32 `mc:"VarInt"`
	SaturationModifier float32
	CanAlwaysEat       bool
}

func (*Food) ID() string {
	return "minecraft:food"
}
