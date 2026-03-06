package component

//codec:gen
type CustomModelData struct {
	Floats  []float32
	Flags   []bool
	Strings []string
	Colors  []int32 `mc:"VarInt"`
}

func (*CustomModelData) ID() string {
	return "minecraft:custom_model_data"
}
