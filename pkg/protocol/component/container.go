package component

import (
	slot2 "git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Container struct {
	Items []slot2.Slot
}

func (*Container) Type() slot2.ComponentID {
	return 66
}

func (*Container) ID() string {
	return "minecraft:container"
}
