package component

import (
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type Equippable struct {
	Slot         int32 `mc:"VarInt"` // 0=mainhand, 1=feet, 2=legs, etc.
	EquipSoundID int32 `mc:"VarInt"`
	//opt:id:EquipSoundID
	EquipSoundEvent *SoundEvent
	HasModel        bool
	//opt:optional:HasModel
	Model            string `mc:"Identifier"`
	HasCameraOverlay bool
	//opt:optional:HasCameraOverlay
	CameraOverlay      string `mc:"Identifier"`
	HasAllowedEntities bool
	//opt:optional:HasAllowedEntities
	AllowedEntitiesID pk.IDSet
	Dispensable       bool
	Swappable         bool
	DamageOnHurt      bool
}

func (*Equippable) Type() slot.ComponentID {
	return 28
}

func (*Equippable) ID() string {
	return "minecraft:equippable"
}
