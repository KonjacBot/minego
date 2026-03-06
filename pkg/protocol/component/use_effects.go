package component

//codec:gen
type UseEffects struct {
	CanSprint          bool
	InteractVibrations bool
	SpeedMultiplier    float32
}

func (u *UseEffects) ID() string {
	return "minecraft:use_effects"
}
