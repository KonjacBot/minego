package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type PotionDurationScale struct {
	EffectMultiplier float32
}

func (*PotionDurationScale) Type() slot.ComponentID {
	return 43
}

func (*PotionDurationScale) ID() string {
	return "minecraft:potion_duration_scale"
}
