package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
	pk "github.com/KonjacBot/go-mc/net/packet"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

type Equippable struct {
	Slot            int32 `mc:"VarInt"` // 0=mainhand, 1=feet, 2=legs, etc.
	EquipSound      packet.OptID[SoundEvent, *SoundEvent]
	AssetID         pk.Option[wire.Identifier, *wire.Identifier]
	CameraOverlay   pk.Option[wire.Identifier, *wire.Identifier]
	AllowedEntities pk.Option[wire.IDSet, *wire.IDSet]
	Dispensable     bool
	Swappable       bool
	DamageOnHurt    bool
	EquipOnInteract bool
	CanBeSheared    bool
	ShearingSound   packet.OptID[SoundEvent, *SoundEvent]
}

func (*Equippable) ID() string {
	return "minecraft:equippable"
}
