package slot

import (
	"bytes"
	"testing"
)

func TestDisplayReadItemUsesCurrentRegistryID(t *testing.T) {
	var display Display
	if _, err := display.ReadFrom(bytes.NewReader([]byte{4, 5})); err != nil {
		t.Fatal(err)
	}

	item, ok := display.SlotDisplay.(*Item)
	if !ok {
		t.Fatalf("display id 4 decoded as %T, want *Item", display.SlotDisplay)
	}
	if item.ID != 5 {
		t.Fatalf("item id = %d, want 5", item.ID)
	}
}
