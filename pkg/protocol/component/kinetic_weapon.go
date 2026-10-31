package component

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type KineticWeapon struct {
	ContactCooldownTicks int32 `mc:"VarInt"`
	DelayTicks           int32 `mc:"VarInt"`
	DismountConditions   pk.Option[KineticWeaponCondition, *KineticWeaponCondition]
	KnockbackConditions  pk.Option[KineticWeaponCondition, *KineticWeaponCondition]
	DamageConditions     pk.Option[KineticWeaponCondition, *KineticWeaponCondition]
	ForwardMovement      float32
	DamageMultiplier     float32
	Sound                pk.Option[pk.OptID[SoundEvent, *SoundEvent], *pk.OptID[SoundEvent, *SoundEvent]]
	HitSound             pk.Option[pk.OptID[SoundEvent, *SoundEvent], *pk.OptID[SoundEvent, *SoundEvent]]
}

//codec:gen
type KineticWeaponCondition struct {
	MaxDurationTicks int32 `mc:"VarInt"`
	MinSpeed         float32
	MinRelativeSpeed float32
}

func (*KineticWeapon) ID() string {
	return "minecraft:kinetic_weapon"
}
