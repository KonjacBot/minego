package client

//codec:gen
type ProjectilePower struct {
	EntityID int32 `mc:"VarInt"`
	Power    float64
}
