package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type LodestoneTracker struct {
	HasGlobalPosition bool
	Dimension         pk.Option[pk.Identifier, *pk.Identifier]
	Position          pk.Option[pk.Position, *pk.Position]
	Tracked           bool
}

func (*LodestoneTracker) ID() string {
	return "minecraft:lodestone_tracker"
}
