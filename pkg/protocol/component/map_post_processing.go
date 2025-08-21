package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type MapPostProcessing struct {
	PostProcessingType int32 `mc:"VarInt"` // 0=Lock, 1=Scale
}

func (*MapPostProcessing) Type() slot.ComponentID {
	return 39
}

func (*MapPostProcessing) ID() string {
	return "minecraft:map_post_processing"
}
