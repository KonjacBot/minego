package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

// Tool 是工具的挖掘規則:26.2 最後多一個 CanDestroyBlocksInCreative。
//
//codec:gen
type Tool struct {
	Rules                      []ToolRule
	DefaultMiningSpeed         float32
	DamagePerBlock             int32 `mc:"VarInt"`
	CanDestroyBlocksInCreative bool
}

// ToolRule 是一條挖掘規則;速度、是否掉落都是 optional(一個 bool 加值),前面沒有另外的 bool。
//
//codec:gen
type ToolRule struct {
	Blocks               pk.IDSet
	Speed                pk.Option[pk.Float, *pk.Float]
	CorrectDropForBlocks pk.Option[pk.Boolean, *pk.Boolean]
}

func (*Tool) ID() string {
	return "minecraft:tool"
}
