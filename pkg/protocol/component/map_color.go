package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type MapColor struct {
	Color int32 `mc:"Int"` // RGB components encoded as integer
}

func (*MapColor) Type() slot.ComponentID {
	return 36
}

func (*MapColor) ID() string {
	return "minecraft:map_color"
}
