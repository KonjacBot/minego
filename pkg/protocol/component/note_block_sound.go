package component

import (
	packet "github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type NoteBlockSound struct {
	Sound packet.Identifier
}

func (*NoteBlockSound) ID() string {
	return "minecraft:note_block_sound"
}
