package component

import (
	"github.com/KonjacBot/go-mc/chat"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type CustomName struct {
	Name chat.Message
}

func (*CustomName) Type() slot.ComponentID {
	return 5
}

func (*CustomName) ID() string {
	return "minecraft:custom_name"
}
