package component

//codec:gen
type BlockState struct {
	Properties []BlockStateProperty
}

//codec:gen
type BlockStateProperty struct {
	Name  string
	Value string
}

func (*BlockState) ID() string {
	return "minecraft:block_state"
}
