package component

import (
	"git.konjactw.dev/falloutBot/go-mc/nbt"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
