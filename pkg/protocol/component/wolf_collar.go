package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type WolfCollar struct {
	Color DyeColor
}

func (*WolfCollar) Type() slot.ComponentID {
	return 75
}

func (*WolfCollar) ID() string {
	return "minecraft:wolf/collar"
}
