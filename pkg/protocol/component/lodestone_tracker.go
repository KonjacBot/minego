package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

//codec:gen
type LodestoneTracker struct {
	HasGlobalPosition bool
	Dimension         pk.Option[packet.Identifier, *packet.Identifier]
	Position          pk.Option[pk.Position, *pk.Position]
	Tracked           bool
}

func (*LodestoneTracker) Type() slot.ComponentID {
	return 58
}

func (*LodestoneTracker) ID() string {
	return "minecraft:lodestone_tracker"
}
