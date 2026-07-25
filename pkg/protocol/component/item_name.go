package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type ItemName struct {
	Name wire.Message
}

func (*ItemName) ID() string {
	return "minecraft:item_name"
}
