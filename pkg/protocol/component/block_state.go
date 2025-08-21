package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type BlockState struct {
	Properties []BlockStateProperty
}

//codec:gen
type BlockStateProperty struct {
	Name  string
	Value string
}

func (*BlockState) Type() slot.ComponentID {
	return 67
}

func (*BlockState) ID() string {
	return "minecraft:block_state"
}
