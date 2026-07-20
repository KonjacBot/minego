package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type GlobalPosition struct {
	Dimension pk.Identifier
	Position  pk.Position
}

type LodestoneTracker struct {
	Target  pk.Option[GlobalPosition, *GlobalPosition]
	Tracked bool
}

func (*LodestoneTracker) ID() string {
	return "minecraft:lodestone_tracker"
}
