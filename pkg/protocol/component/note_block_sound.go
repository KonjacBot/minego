package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type NoteBlockSound struct {
	Sound packet.Identifier
}

func (*NoteBlockSound) Type() slot.ComponentID {
	return 62
}

func (*NoteBlockSound) ID() string {
	return "minecraft:note_block_sound"
}
