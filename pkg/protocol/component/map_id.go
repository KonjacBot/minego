package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type MapID struct {
	MapID int32 `mc:"VarInt"`
}

func (*MapID) Type() slot.ComponentID {
	return 37
}

func (*MapID) ID() string {
	return "minecraft:map_id"
}
