package component

import (
	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type PaintingVariant struct {
	Width   int32
	Height  int32
	AssetID pk.Identifier
	Title   pk.Option[chat.Message, *chat.Message]
	Author  pk.Option[chat.Message, *chat.Message]
}

func (*PaintingVariant) ID() string {
	return "minecraft:painting/variant"
}
