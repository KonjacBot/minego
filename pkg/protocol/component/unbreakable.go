package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
