package component

import (
	"io"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/net/packet"
)

// JukeboxPlayable 是唱片能放的歌:26.2 只有 Holder<JukeboxSong>,
// VarInt 大於 0 是登錄表編號加 1,0 接內嵌的 JukeboxSongData(舊版開頭還有一個 mode byte,已經拿掉)。
type JukeboxPlayable struct {
	SongData packet.OptID[JukeboxSongData, *JukeboxSongData]
}

func (p *JukeboxPlayable) ReadFrom(r io.Reader) (int64, error) {
	return p.SongData.ReadFrom(r)
}

func (p JukeboxPlayable) WriteTo(w io.Writer) (int64, error) {
	return p.SongData.WriteTo(w)
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
