package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type UseCooldown struct {
	Seconds       float32
	CooldownGroup pk.Option[pk.Identifier, *pk.Identifier]
}

func (*UseCooldown) ID() string {
	return "minecraft:use_cooldown"
}
