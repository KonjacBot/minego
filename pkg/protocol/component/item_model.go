package component

import (
	packet "github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type ItemModel struct {
	Model packet.Identifier
}

func (*ItemModel) ID() string {
	return "minecraft:item_model"
}
