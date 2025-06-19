package component

import (
	"git.konjactw.dev/patyhank/minego/codec/slot"
	pk "github.com/Tnze/go-mc/net/packet"
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
