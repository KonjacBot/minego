package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type CustomModelData struct {
	Floats  []float32
	Flags   []bool
	Strings []string
	Colors  []int32 `mc:"VarInt"`
}

func (*CustomModelData) Type() slot.ComponentID {
	return 14
}

func (*CustomModelData) ID() string {
	return "minecraft:custom_model_data"
}
