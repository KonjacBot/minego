package component

//codec:gen
type SuspiciousStewEffects struct {
	Effects []SuspiciousStewEffect
}

//codec:gen
type SuspiciousStewEffect struct {
	TypeID   int32 `mc:"VarInt"`
	Duration int32 `mc:"VarInt"`
}

func (*SuspiciousStewEffects) ID() string {
	return "minecraft:suspicious_stew_effects"
}
