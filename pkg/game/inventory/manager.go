package inventory

import (
	"context"
	"fmt"
	"sync"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

// Manager 管理inventory和container
type Manager struct {
	mu                   sync.RWMutex
	c                    bot.Client
	inventory            *Container
	container            *Container
	cursor               *slot.Slot
	currentContainerID   int32
	currentContainerType int32
}

func NewManager(c bot.Client) *Manager {
	m := &Manager{
		c:                  c,
		inventory:          NewContainerWithSize(c, 0, 46),
		currentContainerID: -1,
	}
	m.inventory.manager = m

	bot.AddHandler(c, func(ctx context.Context, p *client.SetContainerContent) {
		m.mu.Lock()
		matched := false
		if p.WindowID == 0 {
			m.inventory.SetSlots(p.Slots)
			m.inventory.setStateID(p.StateID)
			matched = true
		} else if p.WindowID == m.currentContainerID && m.container != nil {
			m.container.SetSlots(p.Slots)
			m.container.setStateID(p.StateID)
			m.syncPlayerInventory(p.Slots)
			matched = true
		}
		cursor := p.CarriedItem
		if matched {
			m.cursor = &cursor
		}
		m.mu.Unlock()
		if matched {
			m.c.Player().UpdateStateID(p.StateID)
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.BlockChangedAck) {
		m.c.Player().UpdateSequence(p.Sequence)
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.ContainerSetSlot) {
		m.mu.Lock()
		matched := false
		if p.ContainerID == -1 && p.Slot == -1 {
			cursor := p.ItemStack
			m.cursor = &cursor
			matched = true
		} else if p.ContainerID == 0 {
			m.inventory.SetSlot(int(p.Slot), p.ItemStack)
			m.inventory.setStateID(p.StateID)
			matched = true
		} else if p.ContainerID == m.currentContainerID && m.container != nil {
			m.container.SetSlot(int(p.Slot), p.ItemStack)
			m.container.setStateID(p.StateID)
			if count := m.container.SlotCount(); count >= 36 && int(p.Slot) >= count-36 {
				start := count - 36
				m.inventory.SetSlot(9+int(p.Slot)-start, p.ItemStack)
			}
			matched = true
		}
		m.mu.Unlock()
		if matched && !(p.ContainerID == -1 && p.Slot == -1) {
			m.c.Player().UpdateStateID(p.StateID)
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.SetCursorItem) {
		m.mu.Lock()
		cursor := p.CarriedItem
		m.cursor = &cursor
		m.mu.Unlock()
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.SetPlayerInventory) {
		m.mu.Lock()
		if index, ok := playerInventoryMenuSlot(p.Slot); ok {
			m.inventory.SetSlot(index, p.Data)
		}
		if m.container != nil && m.container.SlotCount() >= 36 {
			if offset, ok := playerInventoryWindowOffset(p.Slot); ok {
				m.container.SetSlot(m.container.SlotCount()-36+offset, p.Data)
			}
		}
		m.mu.Unlock()
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.CloseContainer) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if p.WindowID == m.currentContainerID {
			m.currentContainerID = -1
			m.currentContainerType = -1
			m.container = nil
			m.cursor = nil
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.OpenScreen) {
		m.mu.Lock()
		m.currentContainerID = p.WindowID
		m.currentContainerType = p.WindowType
		m.container = NewContainer(c, p.WindowID)
		m.container.manager = m
		m.mu.Unlock()
		_ = bot.PublishEvent(m.c, ContainerOpenEvent{
			WindowID: p.WindowID,
			Type:     p.WindowType,
			Title:    p.WindowTitle,
		})
	})

	return m
}

func (m *Manager) Inventory() bot.Container {
	return m.inventory
}

func (m *Manager) Container() bot.Container {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.container == nil {
		return nil
	}
	return m.container
}
func (m *Manager) Cursor() *slot.Slot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.cursor == nil {
		return nil
	}
	cursor := m.cursor.Clone()
	return &cursor
}

func (m *Manager) CurrentContainerID() int32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentContainerID
}

func (m *Manager) Close() {
	m.mu.Lock()
	id := m.currentContainerID
	m.currentContainerID = -1
	m.currentContainerType = -1
	m.container = nil
	m.cursor = nil
	m.mu.Unlock()
	if id >= 0 {
		_ = m.c.WritePacket(context.Background(), &server.ContainerClose{WindowID: id})
	}
}

// Click 點擊容器slot
func (m *Manager) Click(id int32, slotIndex int16, mode int32, button int32) error {
	return m.click(id, slotIndex, mode, button)
}

func (m *Manager) click(id int32, slotIndex int16, mode int32, button int32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	target := m.inventory
	if id != 0 {
		if id != m.currentContainerID || m.container == nil {
			return fmt.Errorf("container %d is not open", id)
		}
		target = m.container
	}
	cursor := slot.Slot{}
	if m.cursor != nil {
		cursor = m.cursor.Clone()
	}
	layout := ClickLayout{}
	if id > 0 && m.currentContainerType >= 0 && m.currentContainerType <= 5 {
		layout = ClickLayout{PlayerSlotStart: target.SlotCount() - 36, GenericContainerMenu: true}
	}
	packet, prediction, err := BuildClickPacketWithLayout(id, target.StateID(), target.Slots(), cursor, slotIndex, mode, button, layout)
	if err != nil {
		return err
	}
	if err := m.c.WritePacket(context.Background(), packet); err != nil {
		return err
	}
	if prediction.Complete {
		target.SetSlots(prediction.Slots)
		predictedCursor := prediction.Cursor.Clone()
		m.cursor = &predictedCursor
		if id > 0 {
			m.syncPlayerInventory(prediction.Slots)
		}
	}
	return nil
}

func (m *Manager) syncPlayerInventory(windowSlots []slot.Slot) {
	if len(windowSlots) < 36 {
		return
	}
	start := len(windowSlots) - 36
	for index := range 36 {
		m.inventory.SetSlot(9+index, windowSlots[start+index])
	}
}

func playerInventoryMenuSlot(index int32) (int, bool) {
	switch {
	case index >= 0 && index <= 8:
		return 36 + int(index), true
	case index >= 9 && index <= 35:
		return int(index), true
	default:
		return 0, false
	}
}

func playerInventoryWindowOffset(index int32) (int, bool) {
	switch {
	case index >= 0 && index <= 8:
		return 27 + int(index), true
	case index >= 9 && index <= 35:
		return int(index) - 9, true
	default:
		return 0, false
	}
}
