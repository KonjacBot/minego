package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type BlocksAttacks struct {
	BlockDelaySeconds    float32
	DisableCooldownScale float32
	DamageReductions     []DamageReduction
	ItemDamageThreshold  float32
	ItemDamageBase       float32
	ItemDamageFactor     float32
	BypassedBy           pk.Option[pk.Identifier, *pk.Identifier]
	BlockSound           pk.Option[SoundEvent, *SoundEvent]
	DisableSound         pk.Option[SoundEvent, *SoundEvent]
}

//codec:gen
type DamageReduction struct {
	HorizontalBlockingAngle float32
	Type                    pk.Option[pk.IDSet, *pk.IDSet]
	Base                    float32
	Factor                  float32
}

func (*BlocksAttacks) Type() slot.ComponentID {
	return 33
}

func (*BlocksAttacks) ID() string {
	return "minecraft:blocks_attacks"
}
