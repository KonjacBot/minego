package component

//codec:gen
type MapID struct {
	MapID int32 `mc:"VarInt"`
}

func (*MapID) ID() string {
	return "minecraft:map_id"
}
