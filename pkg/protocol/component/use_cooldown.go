package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type UseCooldown struct {
	Seconds       float32
	CooldownGroup pk.Option[pk.Identifier, *pk.Identifier]
}

func (*UseCooldown) Type() slot.ComponentID {
	return 23
}

func (*UseCooldown) ID() string {
	return "minecraft:use_cooldown"
}
