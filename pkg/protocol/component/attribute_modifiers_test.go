package component

import (
	"bytes"
	"testing"

	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

func TestAttributeModifiersConsumesDisplayForEveryEntry(t *testing.T) {
	var wire bytes.Buffer
	if _, err := pk.VarInt(3).WriteTo(&wire); err != nil {
		t.Fatal(err)
	}
	writeEntry := func(attributeID int32, modifierID string, displayType int32, overrideText string) {
		fields := []pk.FieldEncoder{
			pk.VarInt(attributeID),
			pk.Identifier(modifierID),
			pk.Double(2.5),
			pk.VarInt(0),
			pk.VarInt(6),
			pk.VarInt(displayType),
		}
		for _, field := range fields {
			if _, err := field.WriteTo(&wire); err != nil {
				t.Fatal(err)
			}
		}
		if displayType == 2 {
			if _, err := chat.Text(overrideText).WriteTo(&wire); err != nil {
				t.Fatal(err)
			}
		}
	}
	writeEntry(1, "minecraft:test_default", 0, "")
	writeEntry(2, "minecraft:test_hidden", 1, "")
	writeEntry(3, "minecraft:test_override", 2, "Override")
	wireLength := wire.Len()

	var modifiers AttributeModifiers
	read, err := modifiers.ReadFrom(&wire)
	if err != nil {
		t.Fatal(err)
	}
	if read != int64(wireLength) || wire.Len() != 0 {
		t.Fatalf("attribute modifier decoder consumed %d/%d bytes", read, wireLength)
	}
	if len(modifiers.Modifiers) != 3 {
		t.Fatalf("decoded %d modifiers", len(modifiers.Modifiers))
	}
	for i, wantType := range []int32{0, 1, 2} {
		if got := modifiers.Modifiers[i].Display.Type; got != wantType {
			t.Fatalf("modifier[%d] display type = %d, want %d", i, got, wantType)
		}
	}
	if got := modifiers.Modifiers[2].Display.OverrideText.String(); got != "Override" {
		t.Fatalf("override text = %q", got)
	}
}
