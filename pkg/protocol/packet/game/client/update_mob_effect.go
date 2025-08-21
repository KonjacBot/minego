package client

//codec:gen
type EntityEffect struct {
	EntityID  int32 `mc:"VarInt"`
	EffectID  int32 `mc:"VarInt"`
	Amplifier int32 `mc:"VarInt"`
	Duration  int32 `mc:"VarInt"`
	Flags     int8
}
