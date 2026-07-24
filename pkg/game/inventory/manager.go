package inventory

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

// Snapshot is a detached, lock-consistent view of the canonical inventory and
// current menu state. Handled server inventory packets advance both revisions,
// even when their payload is equal to the current state. Local predictions
// advance only Revision.
type Snapshot struct {
	WindowID              int32
	WindowType            int32
	WindowTitle           string
	Ready                 bool
	StateID               int32
	PlayerSlotStart       int
	ContainerSlots        []slot.Slot
	PlayerInventorySlots  []slot.Slot
	Cursor                slot.Slot
	Revision              uint64
	AuthoritativeRevision uint64
	MenuEpoch             uint64
}

// ClickResult describes the local canonical state after a successfully sent
// click. Complete distinguishes supported prediction, including no-ops, from
// click modes whose resulting state must be supplied by the server.
type ClickResult struct {
	Complete              bool
	Changed               bool
	Revision              uint64
	AuthoritativeRevision uint64
}

// CanonicalState is the context-aware inventory state contract implemented by
// Manager. Callers holding bot.InventoryHandler can type-assert this interface
// without requiring legacy handlers to implement the new API.
type CanonicalState interface {
	Snapshot() Snapshot
	Wait(context.Context, uint64, func(Snapshot) bool) (Snapshot, error)
	ClickTransaction(context.Context, int32, int16, int32, int32) (ClickResult, error)
	ClickContext(context.Context, int32, int16, int32, int32) error
}

// Manager 管理inventory和container
type Manager struct {
	mu                    sync.RWMutex
	c                     bot.Client
	inventory             *Container
	container             *Container
	cursor                *slot.Slot
	currentContainerID    int32
	currentContainerType  int32
	currentContainerTitle string
	currentContainerReady bool
	inventoryReady        bool
	revision              uint64
	authoritativeRevision uint64
	menuEpoch             uint64
	changed               chan struct{}
}

var _ CanonicalState = (*Manager)(nil)

func NewManager(c bot.Client) *Manager {
	m := &Manager{
		c:                    c,
		inventory:            NewContainerWithSize(c, 0, 46),
		currentContainerID:   -1,
		currentContainerType: -1,
		changed:              make(chan struct{}),
	}
	m.inventory.manager = m

	bot.AddHandler(c, func(ctx context.Context, p *client.SetContainerContent) {
		m.mu.Lock()
		matched := false
		if p.WindowID == 0 {
			m.inventory.SetSlots(p.Slots)
			m.inventory.setStateID(p.StateID)
			m.inventoryReady = true
			matched = true
		} else if p.WindowID == m.currentContainerID && m.container != nil {
			m.container.SetSlots(p.Slots)
			m.container.setStateID(p.StateID)
			m.currentContainerReady = true
			m.syncPlayerInventory(p.Slots)
			matched = true
		}
		cursor := p.CarriedItem
		if matched {
			m.cursor = &cursor
			m.advanceLocked(true)
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
		if matched {
			m.advanceLocked(true)
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
		m.advanceLocked(true)
		m.mu.Unlock()
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.SetPlayerInventory) {
		m.mu.Lock()
		handled := false
		if index, ok := playerInventoryMenuSlot(p.Slot); ok {
			m.inventory.SetSlot(index, p.Data)
			handled = true
		}
		if m.container != nil && m.container.SlotCount() >= 36 {
			if offset, ok := playerInventoryWindowOffset(p.Slot); ok {
				m.container.SetSlot(m.container.SlotCount()-36+offset, p.Data)
			}
		}
		if handled {
			m.advanceLocked(true)
		}
		m.mu.Unlock()
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.CloseContainer) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if p.WindowID == m.currentContainerID {
			m.currentContainerID = -1
			m.currentContainerType = -1
			m.currentContainerTitle = ""
			m.currentContainerReady = false
			m.container = nil
			m.cursor = nil
			m.advanceLocked(true)
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *client.OpenScreen) {
		m.mu.Lock()
		m.currentContainerID = p.WindowID
		m.currentContainerType = p.WindowType
		m.currentContainerTitle = p.WindowTitle.ClearString()
		m.currentContainerReady = false
		m.menuEpoch++
		m.container = NewContainer(c, p.WindowID)
		m.container.manager = m
		m.advanceLocked(true)
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

// Snapshot returns a detached atomic view. Mutating its slots does not mutate
// Manager state.
func (m *Manager) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.snapshotLocked()
}

// Wait returns the first snapshot newer than after for which predicate is true.
// A nil predicate accepts any newer snapshot. Registration and the initial
// state check share m.mu, so an update cannot be lost between them.
func (m *Manager) Wait(ctx context.Context, after uint64, predicate func(Snapshot) bool) (Snapshot, error) {
	for {
		if err := ctx.Err(); err != nil {
			return Snapshot{}, err
		}
		m.mu.RLock()
		snapshot := m.snapshotLocked()
		changed := m.changed
		m.mu.RUnlock()
		if snapshot.Revision > after && (predicate == nil || predicate(snapshot)) {
			return snapshot, nil
		}

		select {
		case <-ctx.Done():
			return Snapshot{}, ctx.Err()
		case <-changed:
		}
	}
}

func (m *Manager) Close() {
	m.mu.Lock()
	id := m.currentContainerID
	m.currentContainerID = -1
	m.currentContainerType = -1
	m.currentContainerTitle = ""
	m.currentContainerReady = false
	m.container = nil
	m.cursor = nil
	if id >= 0 {
		m.advanceLocked(false)
	}
	m.mu.Unlock()
	if id >= 0 {
		_ = m.c.WritePacket(context.Background(), &server.ContainerClose{WindowID: id})
	}
}

// Click 點擊容器slot
func (m *Manager) Click(id int32, slotIndex int16, mode int32, button int32) error {
	return m.ClickContext(context.Background(), id, slotIndex, mode, button)
}

func (m *Manager) click(id int32, slotIndex int16, mode int32, button int32) error {
	return m.ClickContext(context.Background(), id, slotIndex, mode, button)
}

// ClickContext serializes state read, packet construction, network write, and
// local prediction commit as one transaction.
func (m *Manager) ClickContext(ctx context.Context, id int32, slotIndex int16, mode int32, button int32) error {
	_, err := m.ClickTransaction(ctx, id, slotIndex, mode, button)
	return err
}

// ClickTransaction serializes state read, packet construction, network write,
// prediction commit, and result revision capture as one transaction.
func (m *Manager) ClickTransaction(ctx context.Context, id int32, slotIndex int16, mode int32, button int32) (ClickResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return ClickResult{}, err
	}
	target := m.inventory
	if id != 0 {
		if id != m.currentContainerID || m.container == nil {
			return ClickResult{}, fmt.Errorf("container %d is not open", id)
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
		return ClickResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return ClickResult{}, err
	}
	if err := m.c.WritePacket(ctx, packet); err != nil {
		return ClickResult{}, err
	}
	changed := false
	if prediction.Complete {
		changed = len(prediction.Changed) > 0 || !reflect.DeepEqual(cursor, prediction.Cursor)
		target.SetSlots(prediction.Slots)
		predictedCursor := prediction.Cursor.Clone()
		m.cursor = &predictedCursor
		if id > 0 {
			m.syncPlayerInventory(prediction.Slots)
		}
		if changed {
			m.advanceLocked(false)
		}
	}
	return ClickResult{
		Complete: prediction.Complete, Changed: changed,
		Revision: m.revision, AuthoritativeRevision: m.authoritativeRevision,
	}, nil
}

func (m *Manager) snapshotLocked() Snapshot {
	windowID := int32(0)
	windowType := int32(-1)
	windowTitle := "Inventory"
	ready := m.inventoryReady
	stateID := m.inventory.StateID()
	containerSlots := m.inventory.Slots()
	playerSlotStart := 9
	if m.currentContainerID >= 0 && m.container != nil {
		windowID = m.currentContainerID
		windowType = m.currentContainerType
		windowTitle = m.currentContainerTitle
		ready = m.currentContainerReady
		stateID = m.container.StateID()
		containerSlots = m.container.Slots()
		playerSlotStart = -1
		if len(containerSlots) >= 36 {
			playerSlotStart = len(containerSlots) - 36
		}
	}
	cursor := slot.Slot{}
	if m.cursor != nil {
		cursor = m.cursor.Clone()
	}
	return Snapshot{
		WindowID: windowID, WindowType: windowType, WindowTitle: windowTitle,
		Ready: ready, StateID: stateID, PlayerSlotStart: playerSlotStart,
		ContainerSlots: containerSlots, PlayerInventorySlots: m.inventory.Slots(), Cursor: cursor,
		Revision: m.revision, AuthoritativeRevision: m.authoritativeRevision, MenuEpoch: m.menuEpoch,
	}
}

func (m *Manager) advanceLocked(authoritative bool) {
	m.revision++
	if authoritative {
		m.authoritativeRevision++
	}
	close(m.changed)
	m.changed = make(chan struct{})
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
