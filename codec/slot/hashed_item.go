package slot

//codec:gen
type AddedHashedComponent struct {
	Type     int32 `mc:"VarInt"`
	DataHash int32
}

//codec:gen
type HashedSlot struct {
	HasItem bool
	//opt:optional:HasItem
	ItemID int32 `mc:"VarInt"`
	//opt:optional:HasItem
	ItemCount int32 `mc:"VarInt"`
	//opt:optional:HasItem
	AddComponents AddedHashedComponent
	//opt:optional:HasItem
	RemovedComponents []int32 `mc:"VarInt"`
}
