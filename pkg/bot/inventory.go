package bot

import (
	"git.konjactw.dev/falloutBot/go-mc/level/item"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type Container interface {
	GetSlot(index int) slot.Slot
	Slots() []slot.Slot
	SlotCount() int
	FindEmpty() int16
	FindItem(itemID item.ID) int16
	Click(slot int16, mode int32, button int32) error
}

type InventoryHandler interface {
	Inventory() Container
	Container() Container
	CurrentContainerID() int32
	Click(container int32, slot int16, mode int32, button int32) error
	Close()
}
