package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/nbt"
)

//codec:gen
type DebugStickState struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*DebugStickState) Type() slot.ComponentID {
	return 48
}

func (*DebugStickState) ID() string {
	return "minecraft:debug_stick_state"
}
