package component

import (
	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

// PaintingVariant 是畫的種類:26.2 是 holder,先一個 VarInt,大於 0 是登錄表編號加 1,
// 等於 0 才接著內嵌的 PaintingVariantData。
//
//codec:gen
type PaintingVariant struct {
	Variant pk.OptID[PaintingVariantData, *PaintingVariantData]
}

// PaintingVariantData 是內嵌的畫:寬、高(格數)、材質與標題、作者。
//
//codec:gen
type PaintingVariantData struct {
	Width   int32 `mc:"VarInt"`
	Height  int32 `mc:"VarInt"`
	AssetID pk.Identifier
	Title   pk.Option[chat.Message, *chat.Message]
	Author  pk.Option[chat.Message, *chat.Message]
}

func (*PaintingVariant) ID() string {
	return "minecraft:painting/variant"
}
