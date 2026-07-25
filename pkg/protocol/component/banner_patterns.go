package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type BannerPatterns struct {
	Layers []BannerLayer
}

//codec:gen
type BannerLayer struct {
	Pattern int32 `mc:"VarInt"`
	//opt:enum:Pattern:0
	AssetID wire.Identifier
	//opt:enum:Pattern:0
	TranslationKey string
	Color          DyeColor
}

//codec:gen
type DyeColor struct {
	ColorID int32 `mc:"VarInt"`
}

func (*BannerPatterns) ID() string {
	return "minecraft:banner_patterns"
}
