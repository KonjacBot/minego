package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type Repairable struct {
	Items pk.IDSet
}

func (*Repairable) Type() slot.ComponentID {
	return 29
}

func (*Repairable) ID() string {
	return "minecraft:repairable"
}
