package inventory

import "testing"

func TestContainerRejectsInvalidSizesAndIndexes(t *testing.T) {
	container := NewContainerWithSize(nil, 1, -1)
	if count := container.SlotCount(); count != 0 {
		t.Fatalf("negative-sized container has %d slots", count)
	}
	container.SetSlot(maxContainerSlots, container.GetSlot(0))
	if count := container.SlotCount(); count != 0 {
		t.Fatalf("oversized slot index expanded container to %d slots", count)
	}
}

func TestContainerClickRejectsNilClient(t *testing.T) {
	container := NewContainer(nil, 1)
	if err := container.Click(0, 0, 0); err == nil {
		t.Fatal("Click() accepted a nil client")
	}
}
