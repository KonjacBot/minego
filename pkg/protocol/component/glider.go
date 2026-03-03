package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
