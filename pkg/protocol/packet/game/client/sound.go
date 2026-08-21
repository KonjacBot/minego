package client

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/component"
)

//codec:gen
type SoundEffect struct {
	Sound                                             pk.OptID[component.SoundEvent, *component.SoundEvent]
	SoundCategory                                     int32 `mc:"VarInt"`
	EffectPositionX, EffectPositionY, EffectPositionZ int32
	Volume, Pitch                                     float32
	Seed                                              int64
}
