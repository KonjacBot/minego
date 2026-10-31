package component

// ChickenVariant 是雞的種類:26.2 是 holderRegistry,只有一個登錄表編號(VarInt),不能內嵌。
//
//codec:gen
type ChickenVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*ChickenVariant) ID() string {
	return "minecraft:chicken/variant"
}
