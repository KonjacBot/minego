package component

import (
	"github.com/KonjacBot/go-mc/net/packet"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type Equippable struct {
	Slot            int32 `mc:"VarInt"` // 0=mainhand, 1=feet, 2=legs, etc.
	EquipSound      packet.OptID[SoundEvent, *SoundEvent]
	AssetID         pk.Option[pk.Identifier, *pk.Identifier]
	CameraOverlay   pk.Option[pk.Identifier, *pk.Identifier]
	AllowedEntities pk.Option[pk.IDSet, *pk.IDSet]
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
