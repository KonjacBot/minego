package component

import (
	"github.com/KonjacBot/go-mc/nbt"
)

//codec:gen
type DebugStickState struct {
	Data nbt.RawMessage `mc:"NBT"`
}

func (*DebugStickState) ID() string {
	return "minecraft:debug_stick_state"
}
