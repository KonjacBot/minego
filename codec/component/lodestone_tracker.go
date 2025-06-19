package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	"github.com/Tnze/go-mc/net/packet"
	pk "github.com/Tnze/go-mc/net/packet"
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
