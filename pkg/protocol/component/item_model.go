package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type ItemModel struct {
	Model packet.Identifier
}

func (*ItemModel) ID() string {
	return "minecraft:item_model"
}
