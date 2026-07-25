package component

import (
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type CustomName struct {
	Name wire.Message
}

func (*CustomName) ID() string {
	return "minecraft:custom_name"
}
