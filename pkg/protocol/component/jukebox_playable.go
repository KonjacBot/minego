package component

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type JukeboxPlayable struct {
	Mode     int8
	Name     packet.Identifier
	SongData packet.OptID[JukeboxSongData, *JukeboxSongData]
}

func (p *JukeboxPlayable) ReadFrom(r io.Reader) (n int64, err error) {
	n1, err := (*packet.Byte)(&p.Mode).ReadFrom(r)
	if err != nil {
		return n1, err
	}
	if p.Mode == 0 {
		n2, err := p.Name.ReadFrom(r)
		return n1 + n2, err
	}

	if p.Mode == 1 {
		n2, err := p.SongData.ReadFrom(r)
		return n1 + n2, err
	}

	return n1, err
}

func (p JukeboxPlayable) WriteTo(w io.Writer) (int64, error) {
	n1, err := (*packet.Byte)(&p.Mode).WriteTo(w)
	if err != nil {
		return n1, err
	}

	if p.Mode == 0 {
		n2, err := p.Name.WriteTo(w)
		return n1 + n2, err
	}

	if p.Mode == 1 {
		n2, err := p.SongData.WriteTo(w)
		return n1 + n2, err
	}
	return n1, err
}

//codec:gen
type JukeboxSongData struct {
	SoundEvent  packet.OptID[SoundEvent, *SoundEvent]
	Description chat.Message
	Duration    float32
	Output      int32 `mc:"VarInt"`
}

func (*JukeboxPlayable) Type() slot.ComponentID {
	return 55
}

func (*JukeboxPlayable) ID() string {
	return "minecraft:jukebox_playable"
}
