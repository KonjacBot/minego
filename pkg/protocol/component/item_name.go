package component

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
