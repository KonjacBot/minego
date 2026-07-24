package inventory

import (
	"context"
	"sync"

	"github.com/KonjacBot/go-mc/level/item"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

// Container 代表一個容器
type Container struct {
	mu          sync.RWMutex
	containerID int32
	stateID     int32
	slots       []slot.Slot
	c           bot.Client
	manager     *Manager
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	if index < 0 || index >= len(c.slots) {
		return slot.Slot{}
	}
	return c.slots[index].Clone()
}

func (c *Container) Slots() []slot.Slot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]slot.Slot, len(c.slots))
	for index, value := range c.slots {
		result[index] = value.Clone()
	}
	return result
}

func (c *Container) SlotCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.slots)
}

func (c *Container) FindEmpty() int16 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i, s := range c.slots {
		if s.Count <= 0 {
			return int16(i)
		}
	}
	return -1
}

func (c *Container) FindItem(itemID item.ID) int16 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i, s := range c.slots {
		if s.ItemID == itemID && s.Count > 0 {
			return int16(i)
		}
	}
	return -1
}

func (c *Container) SetSlot(index int, s slot.Slot) {
	if index < 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// 自動擴容
	for len(c.slots) <= index {
		c.slots = append(c.slots, slot.Slot{})
	}
	if index >= 0 && index < len(c.slots) {
		c.slots[index] = s.Clone()
	}
}

func (c *Container) SetSlots(slots []slot.Slot) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.slots = make([]slot.Slot, len(slots))
	for index, value := range slots {
		c.slots[index] = value.Clone()
	}
}

func (c *Container) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.slots = make([]slot.Slot, 0)
}

func (c *Container) setStateID(stateID int32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stateID = stateID
}

func (c *Container) StateID() int32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stateID
}

func (c *Container) Click(idx int16, mode int32, button int32) error {
	if c.manager != nil {
		return c.manager.click(c.containerID, idx, mode, button)
	}
	clickPacket := &server.ContainerClick{
		WindowID: c.containerID,
		StateID:  c.StateID(),
		Slot:     idx,
		Button:   int8(button),
		Mode:     mode,
	}
	return c.c.WritePacket(context.Background(), clickPacket)
}
