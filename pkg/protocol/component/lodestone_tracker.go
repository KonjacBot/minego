package component

import (
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
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
