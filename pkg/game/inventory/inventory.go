package inventory

import (
	"context"

	"git.konjactw.dev/falloutBot/go-mc/level/item"

	"git.konjactw.dev/patyhank/minego/pkg/bot"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

// Container 代表一個容器
type Container struct {
	containerID int32
	slots       []slot.Slot
	c           bot.Client
}

func NewContainer(c bot.Client, cID int32) *Container {
	return &Container{
		c:           c,
		containerID: cID,
		slots:       make([]slot.Slot, 0),
	}
}

func NewContainerWithSize(c bot.Client, cID, size int32) *Container {
	return &Container{
		c:           c,
		containerID: cID,
		slots:       make([]slot.Slot, size),
	}
}

func (c *Container) GetSlot(index int) slot.Slot {
	if index < 0 || index >= len(c.slots) {
		return slot.Slot{}
	}
	return c.slots[index]
}

func (c *Container) Slots() []slot.Slot {
	return c.slots
}

func (c *Container) SlotCount() int {
	return len(c.slots)
}

func (c *Container) FindEmpty() int16 {
	for i, s := range c.slots {
		if s.Count <= 0 {
			return int16(i)
		}
	}
	return -1
}

func (c *Container) FindItem(itemID item.ID) int16 {
	for i, s := range c.slots {
		if s.ItemID == itemID && s.Count > 0 {
			return int16(i)
		}
	}
	return -1
}

func (c *Container) SetSlot(index int, s slot.Slot) {
	// 自動擴容
	for len(c.slots) <= index {
		c.slots = append(c.slots, slot.Slot{})
	}
	if index >= 0 && index < len(c.slots) {
		c.slots[index] = s
	}
}

func (c *Container) SetSlots(slots []slot.Slot) {
	c.slots = make([]slot.Slot, len(slots))
	copy(c.slots, slots)
}

func (c *Container) Clear() {
	c.slots = make([]slot.Slot, 0)
}

func (c *Container) Click(idx int16, mode int32, button int32) error {
	clickPacket := &server.ContainerClick{
		WindowID: c.containerID,
		StateID:  c.c.Player().StateID(),
		Slot:     idx,
		Button:   int8(button),
		Mode:     mode,
	}
	return c.c.WritePacket(context.Background(), clickPacket)
}
