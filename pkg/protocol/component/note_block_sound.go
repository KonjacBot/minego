package component

import (
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
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
