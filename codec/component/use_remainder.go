package component

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type UseRemainder struct {
	Remainder slot.Slot
}

func (*UseRemainder) Type() slot.ComponentID {
	return 22
}

func (*UseRemainder) ID() string {
	return "minecraft:use_remainder"
}
