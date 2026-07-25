package inventory

import (
	"context"
	"sync"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

// Manager 管理inventory和container
type Manager struct {
	mu                 sync.RWMutex
	c                  bot.Client
	inventory          *Container
	container          *Container
	cursor             *slot.Slot
	currentContainerID int32
}

func NewManager(c bot.Client) *Manager {
	m := &Manager{
		c:                  c,
		inventory:          NewContainerWithSize(c, 0, 46),
		currentContainerID: -1,
	}

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
		m.mu.Unlock()
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.CloseContainer) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if p.WindowID == m.currentContainerID {
			m.resetCurrentContainerLocked()
		}
	})
	bot.AddHandler(c, func(context.Context, *client.Login) {
		m.resetCurrentContainerFromLifecycle()
	})
	bot.AddHandler(c, func(context.Context, *client.Respawn) {
		m.resetCurrentContainerFromLifecycle()
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.OpenScreen) {
		m.mu.Lock()
		m.currentContainerID = p.WindowID
		m.container = NewContainer(c, p.WindowID)
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
	m.resetCurrentContainerLocked()
	m.mu.Unlock()
	if id >= 0 {
		_ = m.c.WritePacket(context.Background(), &server.ContainerClose{WindowID: id})
	}
}

func (m *Manager) resetCurrentContainerFromLifecycle() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.currentContainerID >= 0 {
		m.resetCurrentContainerLocked()
	}
}

func (m *Manager) resetCurrentContainerLocked() {
	m.currentContainerID = -1
	m.container = nil
	m.cursor = nil
}

// Click 點擊容器slot
func (m *Manager) Click(id int32, slotIndex int16, mode int32, button int32) error {
	m.mu.RLock()
	container := m.container
	currentID := m.currentContainerID
	m.mu.RUnlock()
	stateID := m.c.Player().StateID()
	if id == 0 {
		stateID = m.inventory.StateID()
	} else if id == currentID && container != nil {
		stateID = container.StateID()
	}
	clickPacket := &server.ContainerClick{
		WindowID: id,
		StateID:  stateID,
		Slot:     slotIndex,
		Button:   int8(button),
		Mode:     mode,
	}
	return m.c.WritePacket(context.Background(), clickPacket)
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
