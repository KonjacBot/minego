package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/net/packet"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type BannerPatterns struct {
	Layers []BannerLayer
}

//codec:gen
type BannerLayer struct {
	PatternType    int32 `mc:"VarInt"`
	AssetID        pk.Option[packet.Identifier, *packet.Identifier]
	TranslationKey pk.Option[pk.String, *pk.String]
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
