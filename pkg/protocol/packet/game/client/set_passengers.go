package client

//codec:gen
type SetPassengers struct {
	EntityID   int32   `mc:"VarInt"`
	Passengers []int32 `mc:"VarInt"`
}
