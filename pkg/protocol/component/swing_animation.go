package component

//codec:gen
type SwingAnimation struct {
	Type     int32 `mc:"VarInt"`
	Duration int32 `mc:"VarInt"`
}

func (*SwingAnimation) ID() string {
	return "minecraft:swing_animation"
}
