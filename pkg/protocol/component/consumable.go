package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Consumable struct {
	ConsumeSeconds float32
	Animation      int32 `mc:"VarInt"` // 0=none, 1=eat, 2=drink, etc.
	SoundID        int32 `mc:"VarInt"`
	//opt:id:SoundID
	SoundEvent          *SoundEvent
	HasConsumeParticles bool
	Effects             []ConsumeEffect
}

//codec:gen
type SoundEvent struct {
	SoundEventID  packet.Identifier
	HasFixedRange bool
	//opt:optional:HasFixedRange
	FixedRange float32
}

//codec:gen
type ConsumeEffect struct {
	Type int32 `mc:"VarInt"`
	//opt:enum:Type:0
	ApplyEffects []PotionEffect
	//opt:enum:Type:0
	ApplyProbability float32
	//opt:enum:Type:1
	RemoveEffects packet.IDSet
	//opt:enum:Type:3
	TeleportRandomlyDiameter float32
	//opt:enum:Type:4
	PlaySound SoundEvent
}

func (*Consumable) Type() slot.ComponentID {
	return 21
}

func (*Consumable) ID() string {
	return "minecraft:consumable"
}
