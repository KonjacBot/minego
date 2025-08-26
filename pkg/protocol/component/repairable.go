package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
