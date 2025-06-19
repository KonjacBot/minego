package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/net/packet"
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
