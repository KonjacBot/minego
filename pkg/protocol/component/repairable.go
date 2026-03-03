package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
