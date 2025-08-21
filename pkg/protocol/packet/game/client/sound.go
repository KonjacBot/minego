package client

import "git.konjactw.dev/patyhank/minego/pkg/protocol/component"

//codec:gen
type SoundEffect struct {
	SoundID int32 `mc:"VarInt"`
	//opt:id:SoundID
	SoundEvent                                        *component.SoundEvent
	EffectPositionX, EffectPositionY, EffectPositionZ int32
	Volume, Pitch                                     float32
	Seed                                              int64
}
