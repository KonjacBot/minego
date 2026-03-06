package component

//codec:gen
type RabbitVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*RabbitVariant) ID() string {
	return "minecraft:rabbit/variant"
}
