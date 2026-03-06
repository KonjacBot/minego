package component

//codec:gen
type SalmonSize struct {
	SizeType int32 `mc:"VarInt"`
}

func (*SalmonSize) ID() string {
	return "minecraft:salmon/size"
}
