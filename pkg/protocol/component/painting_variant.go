package component

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type PaintingVariant struct {
	Width   int32
	Height  int32
	AssetID pk.Identifier
	Title   pk.Option[chat.Message, *chat.Message]
	Author  pk.Option[chat.Message, *chat.Message]
}

func (*PaintingVariant) Type() slot.ComponentID {
	return 89
}

func (*PaintingVariant) ID() string {
	return "minecraft:painting/variant"
}
