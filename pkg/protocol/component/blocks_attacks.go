package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
	HasBlockSound        bool
	//opt:optional:HasBlockSound
	BlockSoundID int32 `mc:"VarInt"`
	//opt:optional:HasBlockSound
	//opt:id:BlockSoundID
	BlockSound SoundEvent

	HasDisableSound bool
	//opt:optional:HasDisableSound
	DisableSoundID int32 `mc:"VarInt"`
	//opt:optional:HasDisableSound
	//opt:id:DisableSoundID
	DisableSound SoundEvent
}

//codec:gen
type DamageReduction struct {
	HorizontalBlockingAngle float32
	HasType                 bool
	//opt:optional:HasType
	Type   pk.IDSet
	Base   float32
	Factor float32
}

func (*BlocksAttacks) Type() slot.ComponentID {
	return 33
}

func (*BlocksAttacks) ID() string {
	return "minecraft:blocks_attacks"
}
