package inventory

import (
	"reflect"
	"testing"

	"github.com/KonjacBot/go-mc/level/item"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestPredictClickPickupAndPlace(t *testing.T) {
	diamond := slot.Slot{ItemID: item.ID(7), Count: 5}
	picked := PredictClick([]slot.Slot{diamond, slot.Slot{}}, slot.Slot{}, 0, 0, 0)
	if !picked.Complete || picked.Slots[0].Count != 0 || picked.Cursor.Count != 5 {
		t.Fatalf("pickup prediction = %#v", picked)
	}
	placed := PredictClick(picked.Slots, picked.Cursor, 1, 0, 0)
	if !placed.Complete || placed.Slots[1].Count != 5 || placed.Cursor.Count != 0 {
		t.Fatalf("place prediction = %#v", placed)
	}
}

func TestPredictClickRightSplitAndHotbarSwap(t *testing.T) {
	value := slot.Slot{ItemID: item.ID(7), Count: 5}
	split := PredictClick([]slot.Slot{value}, slot.Slot{}, 0, 0, 1)
	if split.Slots[0].Count != 2 || split.Cursor.Count != 3 {
		t.Fatalf("right-click split = %#v", split)
	}

	slots := make([]slot.Slot, 46)
	slots[9] = slot.Slot{ItemID: item.ID(7), Count: 1}
	slots[36] = slot.Slot{ItemID: item.ID(8), Count: 2}
	swapped := PredictClick(slots, slot.Slot{}, 9, 2, 0)
	if !swapped.Complete || swapped.Slots[9].ItemID != item.ID(8) || swapped.Slots[36].ItemID != item.ID(7) {
		t.Fatalf("hotbar swap = %#v", swapped)
	}
}

func TestBuildClickPacketKeepsCursorForUnsimulatedMode(t *testing.T) {
	carried := slot.Slot{ItemID: item.ID(7), Count: 2}
	packet, prediction, err := BuildClickPacket(3, 4, make([]slot.Slot, 63), carried, 0, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if prediction.Complete || !packet.CarriedSlot.HasItem || packet.CarriedSlot.ItemCount != 2 || len(packet.ChangedSlots) != 0 {
		t.Fatalf("unsimulated click = packet %#v prediction %#v", packet, prediction)
	}
}

func TestPredictQuickMoveFromContainerMergesThenFillsInReversePlayerOrder(t *testing.T) {
	value := slot.Slot{ItemID: item.ID(7), Count: 10}
	slots := make([]slot.Slot, 63)
	slots[0] = value
	slots[62] = slot.Slot{ItemID: item.ID(7), Count: 60}
	prediction := PredictClickWithLayout(slots, slot.Slot{}, 0, 1, 0, ClickLayout{
		PlayerSlotStart: 27, GenericContainerMenu: true,
	})
	if !prediction.Complete || prediction.Slots[0].Count != 0 || prediction.Slots[62].Count != 64 || prediction.Slots[61].Count != 6 {
		t.Fatalf("container quick move = %#v", prediction)
	}
	wantChanged := []int16{0, 61, 62}
	if !reflect.DeepEqual(prediction.Changed, wantChanged) {
		t.Fatalf("changed slots = %v, want %v", prediction.Changed, wantChanged)
	}
}

func TestPredictQuickMoveFromPlayerUsesForwardContainerOrder(t *testing.T) {
	slots := make([]slot.Slot, 63)
	slots[27] = slot.Slot{ItemID: item.ID(7), Count: 5}
	slots[0] = slot.Slot{ItemID: item.ID(7), Count: 62}
	prediction := PredictClickWithLayout(slots, slot.Slot{}, 27, 1, 0, ClickLayout{
		PlayerSlotStart: 27, GenericContainerMenu: true,
	})
	if !prediction.Complete || prediction.Slots[27].Count != 0 || prediction.Slots[0].Count != 64 || prediction.Slots[1].Count != 3 {
		t.Fatalf("player quick move = %#v", prediction)
	}
}
