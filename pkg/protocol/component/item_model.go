package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type ItemModel struct {
	Model packet.Identifier
}

func (*ItemModel) Type() slot.ComponentID {
	return 7
}

func (*ItemModel) ID() string {
	return "minecraft:item_model"
}
