package component

import (
	"git.konjactw.dev/patyhank/minego/codec/slot"
	"github.com/Tnze/go-mc/chat"
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
