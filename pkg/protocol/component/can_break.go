package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type CanBreak struct {
	BlockPredicates []BlockPredicate
}

func (*CanBreak) Type() slot.ComponentID {
	return 12
}

func (*CanBreak) ID() string {
	return "minecraft:can_break"
}
