package component

import "git.konjactw.dev/patyhank/minego/codec/data/slot"

//codec:gen
type OminousBottleAmplifier struct {
	Amplifier int32 `mc:"VarInt"`
}

func (*OminousBottleAmplifier) Type() slot.ComponentID {
	return 54
}

func (*OminousBottleAmplifier) ID() string {
	return "minecraft:ominous_bottle_amplifier"
}
