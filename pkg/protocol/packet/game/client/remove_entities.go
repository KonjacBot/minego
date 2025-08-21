package client

//codec:gen
type RemoveEntities struct {
	EntityIDs []int32 `mc:"VarInt"`
}
