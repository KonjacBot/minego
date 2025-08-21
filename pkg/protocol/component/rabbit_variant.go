package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type RabbitVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*RabbitVariant) Type() slot.ComponentID {
	return 83
}

func (*RabbitVariant) ID() string {
	return "minecraft:rabbit/variant"
}
