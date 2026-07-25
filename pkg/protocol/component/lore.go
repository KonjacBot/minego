package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type Lore struct {
	Lines []wire.Message
}

func (*Lore) ID() string {
	return "minecraft:lore"
}
