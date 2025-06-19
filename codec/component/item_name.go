package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/chat"
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
