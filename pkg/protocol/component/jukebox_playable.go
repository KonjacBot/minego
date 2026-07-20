package component

import (
	"io"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/net/packet"
)

type JukeboxPlayable struct {
	Song packet.OptID[JukeboxSongData, *JukeboxSongData]
}

func (p *JukeboxPlayable) ReadFrom(r io.Reader) (n int64, err error) {
	*p = JukeboxPlayable{}
	return (&p.Song).ReadFrom(r)
}

func (p JukeboxPlayable) WriteTo(w io.Writer) (int64, error) {
	return (&p.Song).WriteTo(w)
}

//codec:gen
type JukeboxSongData struct {
	SoundEvent  packet.OptID[SoundEvent, *SoundEvent]
	Description chat.Message
	Duration    float32
	Output      int32 `mc:"VarInt"`
}

func (*JukeboxPlayable) ID() string {
	return "minecraft:jukebox_playable"
}
