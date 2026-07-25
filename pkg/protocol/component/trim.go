package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type Trim struct {
	TrimMaterial TrimMaterial
	TrimPattern  TrimPattern
}

//codec:gen
type TrimMaterial struct {
	Suffix      string
	Overrides   []TrimOverride
	Description wire.Message
}

//codec:gen
type TrimOverride struct {
	MaterialType      wire.Identifier
	OverrideAssetName string
}

//codec:gen
type TrimPattern struct {
	AssetName    string
	TemplateItem int32 `mc:"VarInt"`
	Description  wire.Message
	Decal        bool
}

func (*Trim) ID() string {
	return "minecraft:trim"
}
