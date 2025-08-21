package client

//codec:gen
type SetCenterChunk struct {
	X, Z int32 `mc:"VarInt"`
}
