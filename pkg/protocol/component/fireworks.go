package component

//codec:gen
type Fireworks struct {
	FlightDuration int32 `mc:"VarInt"`
	Explosions     []FireworkExplosionData
}

func (*Fireworks) ID() string {
	return "minecraft:fireworks"
}
