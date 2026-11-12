package component

import (
	"github.com/KonjacBot/minego/v2/pkg/protocol/slot"
)

//codec:gen
type SulfurCubeContent struct {
	AbsorbedBlockItemStack slot.ItemStackTemplate
}

func (*SulfurCubeContent) ID() string {
	return "minecraft:sulfur_cube_content"
}
