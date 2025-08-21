package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

//codec:gen
type AttributeModifiers struct {
	Modifiers []AttributeModifier
}

//codec:gen
type AttributeModifier struct {
	AttributeID int32  `mc:"VarInt"`
	ModifierID  string `mc:"Identifier"`
	Value       float64
	Operation   int32 `mc:"VarInt"` // 0=Add, 1=Multiply base, 2=Multiply total
	Slot        int32 `mc:"VarInt"` // 0=Any, 1=Main hand, 2=Offhand, etc.
}

func (*AttributeModifiers) Type() slot.ComponentID {
	return 13
}

func (*AttributeModifiers) ID() string {
	return "minecraft:attribute_modifiers"
}
