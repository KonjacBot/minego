package client

//codec:gen
type RemoveMobEffect struct {
	EntityID int32 `mc:"VarInt"`
	EffectID int32 `mc:"VarInt"`
}
