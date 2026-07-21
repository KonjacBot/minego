package component

//codec:gen
type PotionContents struct {
	HasPotionID bool
	//opt:optional:HasPotionID
	PotionID       int32 `mc:"VarInt"`
	HasCustomColor bool
	//opt:optional:HasCustomColor
	CustomColor   int32 `mc:"Int"`
	CustomEffects []PotionEffect
	HasCustomName bool
	//opt:optional:HasCustomName
	CustomName string
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

func (*PotionContents) ID() string {
	return "minecraft:potion_contents"
}
