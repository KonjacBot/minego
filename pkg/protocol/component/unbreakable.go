package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type Unbreakable struct {
	// no fields
}

func (*Unbreakable) Type() slot.ComponentID {
	return 4
}

func (*Unbreakable) ID() string {
	return "minecraft:unbreakable"
}
