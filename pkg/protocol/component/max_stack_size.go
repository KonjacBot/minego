package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type MaxStackSize struct {
	Size int32 `mc:"VarInt"`
}

func (*MaxStackSize) Type() slot.ComponentID {
	return 1
}

func (*MaxStackSize) ID() string {
	return "minecraft:max_stack_size"
}
