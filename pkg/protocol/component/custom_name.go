package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
	"github.com/Tnze/go-mc/chat"
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
