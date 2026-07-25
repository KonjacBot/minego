package component

import (
	packet "github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type ProvidesBannerPatterns struct {
	Key packet.Identifier
}

func (*ProvidesBannerPatterns) ID() string {
	return "minecraft:provides_banner_patterns"
}
