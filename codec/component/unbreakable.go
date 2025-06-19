package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

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
