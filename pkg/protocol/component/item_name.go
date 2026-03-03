package component

import (
	"github.com/KonjacBot/go-mc/chat"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type ItemName struct {
	Name chat.Message
}

func (*ItemName) Type() slot.ComponentID {
	return 6
}

func (*ItemName) ID() string {
	return "minecraft:item_name"
}
