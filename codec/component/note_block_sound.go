package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/net/packet"
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
