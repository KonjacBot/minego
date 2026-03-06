package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type NoteBlockSound struct {
	Sound packet.Identifier
}

func (*NoteBlockSound) ID() string {
	return "minecraft:note_block_sound"
}
