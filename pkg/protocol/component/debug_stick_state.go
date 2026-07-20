package component

type DebugStickState struct {
	Properties map[string]string
}

func (*DebugStickState) ID() string {
	return "minecraft:debug_stick_state"
}
