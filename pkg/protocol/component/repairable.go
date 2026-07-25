package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type Repairable struct {
	Items wire.IDSet
}

func (*Repairable) ID() string {
	return "minecraft:repairable"
}
