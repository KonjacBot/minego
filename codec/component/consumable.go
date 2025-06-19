package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type Consumable struct {
	ConsumeSeconds      float32
	Animation           int32 `mc:"VarInt"` // 0=none, 1=eat, 2=drink, etc.
	Sound               SoundEvent
	HasConsumeParticles bool
	Effects             []ConsumeEffect
}

//codec:gen
type SoundEvent struct {
	SoundEventID packet.Identifier
	FixedRange   packet.Option[packet.Float, *packet.Float]
}

//codec:gen
type ConsumeEffect struct {
	Type int32 `mc:"VarInt"`
	// Data varies by type - would need custom handling
}

func (*Consumable) Type() slot.ComponentID {
	return 21
}

func (*Consumable) ID() string {
	return "minecraft:consumable"
}
