package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
