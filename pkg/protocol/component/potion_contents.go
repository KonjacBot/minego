package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type PotionContents struct {
	PotionID      pk.Option[pk.VarInt, *pk.VarInt]
	CustomColor   pk.Option[pk.Int, *pk.Int]
	CustomEffects []PotionEffect
	CustomName    pk.Option[pk.String, *pk.String]
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
	HiddenEffect    *PotionEffect
}

func (*PotionContents) ID() string {
	return "minecraft:potion_contents"
}
