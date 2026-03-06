package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type LodestoneTracker struct {
	HasGlobalPosition bool
	Dimension         pk.Option[packet.Identifier, *packet.Identifier]
	Position          pk.Option[pk.Position, *pk.Position]
	Tracked           bool
}

func (*LodestoneTracker) ID() string {
	return "minecraft:lodestone_tracker"
}
