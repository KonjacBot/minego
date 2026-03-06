package component

//codec:gen
type MapPostProcessing struct {
	PostProcessingType int32 `mc:"VarInt"` // 0=Lock, 1=Scale
}

func (*MapPostProcessing) ID() string {
	return "minecraft:map_post_processing"
}
