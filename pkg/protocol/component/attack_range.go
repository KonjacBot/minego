package component

//codec:gen
type AttackRange struct {
	MinRange float32
	MaxRange float32

	MinCreativeRange float32
	MaxCreativeRange float32

	HitboxMargin float32

	MobFactor float32
}

func (*AttackRange) ID() string {
	return "minecraft:attack_range"
}
