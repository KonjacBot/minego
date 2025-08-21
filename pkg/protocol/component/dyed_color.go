package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type DyedColor struct {
	Color int32 `mc:"Int"` // RGB components encoded as integer
}

func (*DyedColor) Type() slot.ComponentID {
	return 35
}

func (*DyedColor) ID() string {
	return "minecraft:dyed_color"
}
