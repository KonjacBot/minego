package component

//codec:gen
type MinimumAttackCharge struct {
	Charge float32
}

func (c *MinimumAttackCharge) ID() string {
	return "minecraft:minimum_attack_charge"
}
