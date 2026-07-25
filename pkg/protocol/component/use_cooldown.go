package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type UseCooldown struct {
	Seconds       float32
	CooldownGroup pk.Option[wire.Identifier, *wire.Identifier]
}

func (*UseCooldown) ID() string {
	return "minecraft:use_cooldown"
}
