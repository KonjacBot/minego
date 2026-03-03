package component

import (
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
