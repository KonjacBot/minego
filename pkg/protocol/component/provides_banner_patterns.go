package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type ProvidesBannerPatterns struct {
	Key packet.Identifier
}

func (*ProvidesBannerPatterns) ID() string {
	return "minecraft:provides_banner_patterns"
}
