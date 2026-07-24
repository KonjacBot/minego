package component

import "github.com/KonjacBot/go-mc/chat"

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
	Display     AttributeModifierDisplay
}

//codec:gen
type AttributeModifierDisplay struct {
	Type int32 `mc:"VarInt"` // 0=Default, 1=Hidden, 2=Override text
	//opt:enum:Type:2
	OverrideText chat.Message
}

func (*AttributeModifiers) ID() string {
	return "minecraft:attribute_modifiers"
}
