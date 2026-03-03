package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
