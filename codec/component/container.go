package component

import "git.konjactw.dev/patyhank/minego/codec/slot"

//codec:gen
type Container struct {
	Items []slot.Slot
}

func (*Container) Type() slot.ComponentID {
	return 66
}

func (*Container) ID() string {
	return "minecraft:container"
}
