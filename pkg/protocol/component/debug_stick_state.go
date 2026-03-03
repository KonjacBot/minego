package component

import (
	"github.com/KonjacBot/go-mc/nbt"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
