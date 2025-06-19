package component

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type ChargedProjectiles struct {
	Projectiles []slot.Slot
}

func (*ChargedProjectiles) Type() slot.ComponentID {
	return 40
}

func (*ChargedProjectiles) ID() string {
	return "minecraft:charged_projectiles"
}
