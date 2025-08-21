package client

//codec:gen
type UpdateEntityPosition struct {
	EntityID               int32 `mc:"VarInt"`
	DeltaX, DeltaY, DeltaZ int16
	OnGround               bool
}
