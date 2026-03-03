package component

import (
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type ProvidesBannerPatterns struct {
	Key packet.Identifier
}

func (*ProvidesBannerPatterns) Type() slot.ComponentID {
	return 56
}

func (*ProvidesBannerPatterns) ID() string {
	return "minecraft:provides_banner_patterns"
}
