package component

import "github.com/KonjacBot/go-mc/level/item"

type PotDecorations struct {
	Back  *item.ID
	Left  *item.ID
	Right *item.ID
	Front *item.ID
}

func (*PotDecorations) ID() string {
	return "minecraft:pot_decorations"
}
