package component

//codec:gen
type OminousBottleAmplifier struct {
	Amplifier int32 `mc:"VarInt"`
}

func (*OminousBottleAmplifier) ID() string {
	return "minecraft:ominous_bottle_amplifier"
}
