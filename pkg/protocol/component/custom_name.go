package component

import (
	"github.com/KonjacBot/go-mc/chat"
)

//codec:gen
type CustomName struct {
	Name chat.Message
}

func (*CustomName) ID() string {
	return "minecraft:custom_name"
}
