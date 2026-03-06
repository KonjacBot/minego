package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type Tool struct {
	Rules              []ToolRule
	DefaultMiningSpeed float32
	DamagePerBlock     int32 `mc:"VarInt"`
}

//codec:gen
type ToolRule struct {
	Blocks                  pk.IDSet
	HasSpeed                bool
	Speed                   pk.Option[pk.Float, *pk.Float]
	HasCorrectDropForBlocks bool
	CorrectDropForBlocks    pk.Option[pk.Boolean, *pk.Boolean]
}

func (*Tool) ID() string {
	return "minecraft:tool"
}
