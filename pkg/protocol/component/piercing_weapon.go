package component

import pk "github.com/KonjacBot/go-mc/net/packet"

//codec:gen
type PiercingWeapon struct {
	DealsKnockback bool
	Dismounts      bool
	Sound          pk.Option[SoundEvent, *SoundEvent]
	HitSound       pk.Option[SoundEvent, *SoundEvent]
}

func (*PiercingWeapon) ID() string {
	return "minecraft:piercing_weapon"
}
