package component

import (
	pk "github.com/KonjacBot/go-mc/net/packet"
)

// Equippable 是能穿戴的物品:音效是 holder(VarInt 大於 0 是編號加 1,0 接內嵌),
// 26.2 在 DamageOnHurt 之後多了右鍵穿上、能不能剪、剪的音效。
//
//codec:gen
type Equippable struct {
	Slot            int32 `mc:"VarInt"` // 0=mainhand, 1=feet, 2=legs, etc.
	EquipSound      pk.OptID[SoundEvent, *SoundEvent]
	Model           pk.Option[pk.Identifier, *pk.Identifier]
	CameraOverlay   pk.Option[pk.Identifier, *pk.Identifier]
	AllowedEntities pk.Option[pk.IDSet, *pk.IDSet]
	Dispensable     bool
	Swappable       bool
	DamageOnHurt    bool
	EquipOnInteract bool
	CanBeSheared    bool
	ShearingSound   pk.OptID[SoundEvent, *SoundEvent]
}

func (*Equippable) ID() string {
	return "minecraft:equippable"
}
