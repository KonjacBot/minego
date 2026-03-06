package component

import (
	"github.com/KonjacBot/go-mc/chat"
)

//codec:gen
type Lore struct {
	Lines []chat.Message
}

func (*Lore) ID() string {
	return "minecraft:lore"
}
