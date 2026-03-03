package component

import (
	"github.com/KonjacBot/go-mc/chat"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type Lore struct {
	Lines []chat.Message
}

func (*Lore) Type() slot.ComponentID {
	return 8
}

func (*Lore) ID() string {
	return "minecraft:lore"
}
