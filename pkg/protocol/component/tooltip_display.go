package component

//codec:gen
type TooltipDisplay struct {
	HideTooltip      bool
	HiddenComponents []int32 `mc:"VarInt"`
}

func (*TooltipDisplay) ID() string {
	return "minecraft:tooltip_display"
}
