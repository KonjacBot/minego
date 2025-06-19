package client

//codec:gen
type OpenHorseScreen struct {
	WindowID             int32 `mc:"VarInt"`
	InventoryColumnsSize int32 `mc:"VarInt"`
	EntityID             int32
}
