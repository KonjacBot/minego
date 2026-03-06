package component

//codec:gen
type LlamaVariant struct {
	Variant int32 `mc:"VarInt"`
}

func (*LlamaVariant) ID() string {
	return "minecraft:llama/variant"
}
