package component

import (
	"git.konjactw.dev/patyhank/minego/codec/data/slot"
	pk "github.com/Tnze/go-mc/net/packet"
)

//codec:gen
type Equippable struct {
	Slot               int32 `mc:"VarInt"` // 0=mainhand, 1=feet, 2=legs, etc.
	EquipSound         SoundEvent
	HasModel           bool
	Model              pk.Option[pk.Identifier, *pk.Identifier]
	HasCameraOverlay   bool
	CameraOverlay      pk.Option[pk.Identifier, *pk.Identifier]
	HasAllowedEntities bool
	AllowedEntities    pk.Option[pk.IDSet, *pk.IDSet]
	Dispensable        bool
	Swappable          bool
	DamageOnHurt       bool
}

func (*Equippable) Type() slot.ComponentID {
	return 28
}

func (*Equippable) ID() string {
	return "minecraft:equippable"
}
