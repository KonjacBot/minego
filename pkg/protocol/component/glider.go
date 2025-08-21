package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Glider struct {
	// no fields
}

func (*Glider) Type() slot.ComponentID {
	return 30
}

func (*Glider) ID() string {
	return "minecraft:glider"
}
