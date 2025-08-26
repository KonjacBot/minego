package component

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
