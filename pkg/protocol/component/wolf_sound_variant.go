package component

//codec:gen
type WolfSoundVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*WolfSoundVariant) ID() string {
	return "minecraft:wolf/sound_variant"
}
