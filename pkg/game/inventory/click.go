package inventory

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

const defaultMaxStackSize int32 = 64

// ClickPrediction contains the client-side menu state after a click. Complete
// is false for menu actions whose destination depends on server-owned layout
// rules, but CarriedSlot still reflects the current authoritative cursor.
type ClickPrediction struct {
	Slots    []slot.Slot
	Cursor   slot.Slot
	Changed  []int16
	Complete bool
}

type ClickLayout struct {
	PlayerSlotStart      int
	GenericContainerMenu bool
}

func BuildClickPacket(windowID, stateID int32, slots []slot.Slot, cursor slot.Slot, slotIndex int16, mode, button int32) (*server.ContainerClick, ClickPrediction, error) {
	return BuildClickPacketWithLayout(windowID, stateID, slots, cursor, slotIndex, mode, button, ClickLayout{})
}

func BuildClickPacketWithLayout(windowID, stateID int32, slots []slot.Slot, cursor slot.Slot, slotIndex int16, mode, button int32, layout ClickLayout) (*server.ContainerClick, ClickPrediction, error) {
	prediction := PredictClickWithLayout(slots, cursor, slotIndex, mode, button, layout)
	carried, err := slot.HashSlot(prediction.Cursor)
	if err != nil {
		return nil, ClickPrediction{}, fmt.Errorf("hash carried slot: %w", err)
	}
	packet := &server.ContainerClick{
		WindowID:    windowID,
		StateID:     stateID,
		Slot:        slotIndex,
		Button:      int8(button),
		Mode:        mode,
		CarriedSlot: carried,
	}
	for _, index := range prediction.Changed {
		hashed, err := slot.HashSlot(prediction.Slots[index])
		if err != nil {
			return nil, ClickPrediction{}, fmt.Errorf("hash changed slot %d: %w", index, err)
		}
		packet.ChangedSlots = append(packet.ChangedSlots, server.ChangedSlot{Slot: index, SlotData: hashed})
	}
	return packet, prediction, nil
}

func PredictClick(slots []slot.Slot, cursor slot.Slot, slotIndex int16, mode, button int32) ClickPrediction {
	return PredictClickWithLayout(slots, cursor, slotIndex, mode, button, ClickLayout{})
}

func PredictClickWithLayout(slots []slot.Slot, cursor slot.Slot, slotIndex int16, mode, button int32, layout ClickLayout) ClickPrediction {
	prediction := ClickPrediction{
		Slots:  append([]slot.Slot(nil), slots...),
		Cursor: cursor.Clone(),
	}
	changed := make(map[int16]struct{})
	setSlot := func(index int16, value slot.Slot) {
		if index < 0 || int(index) >= len(prediction.Slots) || reflect.DeepEqual(prediction.Slots[index], value) {
			return
		}
		prediction.Slots[index] = value.Clone()
		changed[index] = struct{}{}
	}

	switch mode {
	case 0:
		prediction.Complete = predictPickup(&prediction, setSlot, slotIndex, button)
	case 1:
		prediction.Complete = layout.GenericContainerMenu && predictQuickMove(&prediction, setSlot, slotIndex, layout.PlayerSlotStart)
	case 2:
		prediction.Complete = predictSwap(&prediction, setSlot, slotIndex, button)
	case 4:
		prediction.Complete = predictThrow(&prediction, setSlot, slotIndex, button)
	default:
		prediction.Complete = false
	}
	for index := range changed {
		prediction.Changed = append(prediction.Changed, index)
	}
	sort.Slice(prediction.Changed, func(i, j int) bool { return prediction.Changed[i] < prediction.Changed[j] })
	return prediction
}

func predictQuickMove(prediction *ClickPrediction, setSlot func(int16, slot.Slot), index int16, playerStart int) bool {
	if index < 0 || int(index) >= len(prediction.Slots) || playerStart <= 0 || playerStart+36 != len(prediction.Slots) {
		return false
	}
	source := prediction.Slots[index].Clone()
	if source.Count <= 0 {
		return true
	}
	start, end, reverse := 0, playerStart, false
	if int(index) < playerStart {
		start, end, reverse = playerStart, len(prediction.Slots), true
	}
	remaining := source.Clone()
	forEachSlot(start, end, reverse, func(targetIndex int) bool {
		target := prediction.Slots[targetIndex].Clone()
		if !stackable(target, remaining) || target.Count >= defaultMaxStackSize {
			return false
		}
		moved := min(defaultMaxStackSize-target.Count, remaining.Count)
		target.Count += moved
		remaining.Count -= moved
		setSlot(int16(targetIndex), target)
		return remaining.Count <= 0
	})
	if remaining.Count > 0 {
		forEachSlot(start, end, reverse, func(targetIndex int) bool {
			if prediction.Slots[targetIndex].Count > 0 {
				return false
			}
			placed := remaining.Clone()
			placed.Count = min(defaultMaxStackSize, remaining.Count)
			remaining.Count -= placed.Count
			setSlot(int16(targetIndex), placed)
			return remaining.Count <= 0
		})
	}
	setSlot(index, normalizeSlot(remaining))
	return true
}

func forEachSlot(start, end int, reverse bool, visit func(int) bool) {
	if reverse {
		for index := end - 1; index >= start; index-- {
			if visit(index) {
				return
			}
		}
		return
	}
	for index := start; index < end; index++ {
		if visit(index) {
			return
		}
	}
}

func predictPickup(prediction *ClickPrediction, setSlot func(int16, slot.Slot), index int16, button int32) bool {
	if button != 0 && button != 1 {
		return false
	}
	if index == -999 {
		if prediction.Cursor.Count <= 0 {
			return true
		}
		if button == 0 || prediction.Cursor.Count == 1 {
			prediction.Cursor = slot.Slot{}
		} else {
			prediction.Cursor.Count--
		}
		return true
	}
	if index < 0 || int(index) >= len(prediction.Slots) {
		return false
	}
	clicked := prediction.Slots[index].Clone()
	carried := prediction.Cursor.Clone()
	if button == 0 {
		switch {
		case carried.Count <= 0:
			prediction.Cursor = clicked
			setSlot(index, slot.Slot{})
		case clicked.Count <= 0:
			setSlot(index, carried)
			prediction.Cursor = slot.Slot{}
		case stackable(clicked, carried):
			moved := min(defaultMaxStackSize-clicked.Count, carried.Count)
			if moved > 0 {
				clicked.Count += moved
				carried.Count -= moved
				setSlot(index, clicked)
				prediction.Cursor = normalizeSlot(carried)
			}
		default:
			setSlot(index, carried)
			prediction.Cursor = clicked
		}
		return true
	}
	switch {
	case carried.Count <= 0 && clicked.Count > 0:
		amount := (clicked.Count + 1) / 2
		prediction.Cursor = clicked
		prediction.Cursor.Count = amount
		clicked.Count -= amount
		setSlot(index, normalizeSlot(clicked))
	case carried.Count > 0 && clicked.Count <= 0:
		placed := carried
		placed.Count = 1
		setSlot(index, placed)
		carried.Count--
		prediction.Cursor = normalizeSlot(carried)
	case carried.Count > 0 && stackable(clicked, carried) && clicked.Count < defaultMaxStackSize:
		clicked.Count++
		carried.Count--
		setSlot(index, clicked)
		prediction.Cursor = normalizeSlot(carried)
	case carried.Count > 0:
		setSlot(index, carried)
		prediction.Cursor = clicked
	}
	return true
}

func predictSwap(prediction *ClickPrediction, setSlot func(int16, slot.Slot), index int16, button int32) bool {
	if index < 0 || int(index) >= len(prediction.Slots) {
		return false
	}
	var swapIndex int
	switch {
	case button >= 0 && button <= 8:
		if len(prediction.Slots) <= 46 {
			swapIndex = 36 + int(button)
		} else {
			swapIndex = len(prediction.Slots) - 9 + int(button)
		}
	case button == 40 && len(prediction.Slots) >= 46:
		swapIndex = 45
	default:
		return false
	}
	if swapIndex < 0 || swapIndex >= len(prediction.Slots) || swapIndex == int(index) {
		return swapIndex == int(index)
	}
	clicked := prediction.Slots[index]
	swapped := prediction.Slots[swapIndex]
	setSlot(index, swapped)
	setSlot(int16(swapIndex), clicked)
	return true
}

func predictThrow(prediction *ClickPrediction, setSlot func(int16, slot.Slot), index int16, button int32) bool {
	if index < 0 || int(index) >= len(prediction.Slots) || (button != 0 && button != 1) {
		return false
	}
	clicked := prediction.Slots[index].Clone()
	if clicked.Count <= 0 {
		return true
	}
	if button == 1 || clicked.Count == 1 {
		setSlot(index, slot.Slot{})
	} else {
		clicked.Count--
		setSlot(index, clicked)
	}
	return true
}

func stackable(left, right slot.Slot) bool {
	return left.Count > 0 && right.Count > 0 && left.ItemID == right.ItemID &&
		reflect.DeepEqual(left.AddComponent, right.AddComponent) &&
		reflect.DeepEqual(left.RemoveComponent, right.RemoveComponent)
}

func normalizeSlot(value slot.Slot) slot.Slot {
	if value.Count <= 0 {
		return slot.Slot{}
	}
	return value
}
