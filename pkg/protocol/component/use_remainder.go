package component

import (
	slot2 "github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type UseRemainder struct {
	Remainder slot2.Slot
}

func (*UseRemainder) Type() slot2.ComponentID {
	return 22
}

func (*UseRemainder) ID() string {
	return "minecraft:use_remainder"
}
