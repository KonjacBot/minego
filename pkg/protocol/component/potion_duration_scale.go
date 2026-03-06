package component

//codec:gen
type PotionDurationScale struct {
	EffectMultiplier float32
}

func (*PotionDurationScale) ID() string {
	return "minecraft:potion_duration_scale"
}
