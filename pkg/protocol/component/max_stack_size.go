package component

//codec:gen
type MaxStackSize struct {
	Size int32 `mc:"VarInt"`
}

func (*MaxStackSize) ID() string {
	return "minecraft:max_stack_size"
}
