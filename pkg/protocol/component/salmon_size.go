package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type SalmonSize struct {
	SizeType int32 `mc:"VarInt"`
}

func (*SalmonSize) Type() slot.ComponentID {
	return 77
}

func (*SalmonSize) ID() string {
	return "minecraft:salmon/size"
}
