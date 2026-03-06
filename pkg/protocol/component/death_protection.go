package component

//codec:gen
type DeathProtection struct {
	Effects []ConsumeEffect
}

func (*DeathProtection) ID() string {
	return "minecraft:death_protection"
}
