package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type BannerPatterns struct {
	Layers []BannerLayer
}

//codec:gen
type BannerLayer struct {
	Pattern int32 `mc:"VarInt"`
	//opt:enum:Pattern:0
	AssetID pk.Identifier
	//opt:enum:Pattern:0
	TranslationKey string
	Color          DyeColor
}

//codec:gen
type DyeColor struct {
	ColorID int32 `mc:"VarInt"`
}

func (*BannerPatterns) Type() slot.ComponentID {
	return 63
}

func (*BannerPatterns) ID() string {
	return "minecraft:banner_patterns"
}
