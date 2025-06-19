package client

//codec:gen
type PickupItem struct {
	CollectedEntityID int32 `mc:"VarInt"`
	CollectorEntityID int32 `mc:"VarInt"`
	PickupItemCount   int32 `mc:"VarInt"`
}
