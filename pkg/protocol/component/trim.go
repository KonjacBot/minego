package component

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
	Description chat.Message
}

//codec:gen
type TrimOverride struct {
	MaterialType      pk.Identifier
	OverrideAssetName string
}

//codec:gen
type TrimPattern struct {
	AssetName    string
	TemplateItem int32 `mc:"VarInt"`
	Description  chat.Message
	Decal        bool
}

func (*Trim) Type() slot.ComponentID {
	return 47
}

func (*Trim) ID() string {
	return "minecraft:trim"
}
