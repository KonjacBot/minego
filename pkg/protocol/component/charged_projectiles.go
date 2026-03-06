package component

import (
	slot2 "github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type ChargedProjectiles struct {
	Projectiles []slot2.Slot
}

func (*ChargedProjectiles) ID() string {
	return "minecraft:charged_projectiles"
}
