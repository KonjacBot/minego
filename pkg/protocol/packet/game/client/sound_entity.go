package client

import (
	"git.konjactw.dev/patyhank/minego/codec/component"
	"github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type EntitySoundEffect struct {
	SoundEvent    packet.OptID[component.SoundEvent, *component.SoundEvent]
	SoundCategory int32 `mc:"VarInt"`
	EntityID      int32 `mc:"VarInt"`
	Volume        float32
	Pitch         float32
	Seed          int32
}
