package component

import (
	"github.com/KonjacBot/go-mc/chat"
)

//codec:gen
type ItemName struct {
	Name chat.Message
}

func (*ItemName) ID() string {
	return "minecraft:item_name"
}
