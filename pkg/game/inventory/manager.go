package inventory

import (
	"context"

	"git.konjactw.dev/patyhank/minego/pkg/bot"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

// Manager 管理inventory和container
type Manager struct {
	c                  bot.Client
	inventory          *Container
	container          *Container
	cursor             *slot.Slot
	currentContainerID int32
}

func NewManager(c bot.Client) *Manager {
	m := &Manager{
		c:                  c,
		inventory:          NewContainerWithSize(c, 0, 45),
		currentContainerID: 0,
	}

	bot.AddHandler(c, func(ctx context.Context, p *client.SetContainerContent) {
		if p.WindowID == 0 {
			m.inventory.SetSlots(p.Slots)
		} else if m.container != nil {
			m.container.SetSlots(p.Slots)
		}
		m.c.Player().UpdateStateID(p.StateID)
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.ContainerSetSlot) {
		if p.ContainerID == 0 {
			m.inventory.SetSlot(int(p.Slot), p.ItemStack)
		} else if m.container != nil {
			m.container.SetSlot(int(p.Slot), p.ItemStack)
		}
		m.c.Player().UpdateStateID(p.StateID)
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.CloseContainer) {
		if p.WindowID == m.currentContainerID {
			m.currentContainerID = -1
			if m.container != nil {
				m.container = nil
			}
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.OpenScreen) {
		m.currentContainerID = p.WindowID
		m.container = NewContainer(c, p.WindowID)
		go bot.PublishEvent(m.c, ContainerOpenEvent{
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
	return m.container
}
func (m *Manager) Cursor() *slot.Slot {
	return m.cursor
}

func (m *Manager) CurrentContainerID() int32 {
	return m.currentContainerID
}

func (m *Manager) Close() {
	if m.currentContainerID != -1 {
		_ = m.c.WritePacket(context.Background(), &server.ContainerClose{WindowID: m.currentContainerID})
		m.currentContainerID = -1
	} else {
		_ = m.c.WritePacket(context.Background(), &server.ContainerClose{WindowID: 0})
		m.currentContainerID = -1
	}
}

// Click 點擊容器slot
func (m *Manager) Click(id int32, slotIndex int16, mode int32, button int32) error {
	clickPacket := &server.ContainerClick{
		WindowID: id,
		StateID:  m.c.Player().StateID(),
		Slot:     slotIndex,
		Button:   int8(button),
		Mode:     mode,
	}
	return m.c.WritePacket(context.Background(), clickPacket)
}
