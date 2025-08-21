package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Rarity struct {
	Rarity int32 `mc:"VarInt"` // 0=Common, 1=Uncommon, 2=Rare, 3=Epic
}

func (*Rarity) Type() slot.ComponentID {
	return 9
}

func (*Rarity) ID() string {
	return "minecraft:rarity"
}
