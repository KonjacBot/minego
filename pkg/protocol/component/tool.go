package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type Tool struct {
	Rules                      []ToolRule
	DefaultMiningSpeed         float32
	DamagePerBlock             int32 `mc:"VarInt"`
	CanDestroyBlocksInCreative bool
}

type ToolRule struct {
	Blocks               pk.IDSet
	Speed                pk.Option[pk.Float, *pk.Float]
	CorrectDropForBlocks pk.Option[pk.Boolean, *pk.Boolean]
}

func (*Tool) ID() string {
	return "minecraft:tool"
}
