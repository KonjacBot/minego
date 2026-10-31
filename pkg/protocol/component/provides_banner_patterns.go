package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
)

// ProvidesBannerPatterns 是旗幟圖案物品能做的圖案:26.2 是 holderSet,
// 開頭 VarInt 0 接標籤名稱,n 大於 0 接 n-1 個登錄表編號。
//
//codec:gen
type ProvidesBannerPatterns struct {
	Patterns packet.IDSet
}

func (*ProvidesBannerPatterns) ID() string {
	return "minecraft:provides_banner_patterns"
}
