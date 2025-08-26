package component

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
