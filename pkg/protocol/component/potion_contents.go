package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type PotionContents struct {
	HasPotionID    bool
	PotionID       int32 `mc:"VarInt"`
	HasCustomColor bool
	CustomColor    int32 `mc:"Int"`
	CustomEffects  []PotionEffect
	CustomName     string
}

//codec:gen
type PotionEffect struct {
	TypeID int32 `mc:"VarInt"`

	Details PotionEffectDetails
}

//codec:gen
type PotionEffectDetails struct {
	Amplifier       int32 `mc:"VarInt"`
	Duration        int32 `mc:"VarInt"`
	Ambient         bool
	ShowParticles   bool
	ShowIcon        bool
	HasHiddenEffect bool
	//opt:optional:HasHiddenEffect
	HiddenEffect *PotionEffect
}

func (*PotionContents) Type() slot.ComponentID {
	return 42
}

func (*PotionContents) ID() string {
	return "minecraft:potion_contents"
}
