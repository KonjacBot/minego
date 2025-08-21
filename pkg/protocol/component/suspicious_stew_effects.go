package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type SuspiciousStewEffects struct {
	Effects []SuspiciousStewEffect
}

//codec:gen
type SuspiciousStewEffect struct {
	TypeID   int32 `mc:"VarInt"`
	Duration int32 `mc:"VarInt"`
}

func (*SuspiciousStewEffects) Type() slot.ComponentID {
	return 44
}

func (*SuspiciousStewEffects) ID() string {
	return "minecraft:suspicious_stew_effects"
}
