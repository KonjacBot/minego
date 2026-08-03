package inventory

import (
	"context"
	"testing"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/packetid"
	packet "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
	gameclient "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestManagerIgnoresContentForStaleWindowAndTracksCursor(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	initial := make([]slot.Slot, 46)
	initial[9] = slot.Slot{ItemID: 12, Count: 1}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, Slots: initial,
	})

	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
	stalePlayer := make([]slot.Slot, 46)
	stalePlayer[9] = slot.Slot{ItemID: 13, Count: 2}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, Slots: stalePlayer,
	})
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 4, StateID: 3, Slots: []slot.Slot{{Count: 1}}, CarriedItem: slot.Slot{Count: 9},
	})
	if got := m.Container().SlotCount(); got != 0 {
		t.Fatalf("stale window changed slot count to %d", got)
	}
	if cursor := m.Cursor(); cursor != nil && cursor.Count > 0 {
		t.Fatalf("stale window changed cursor to %#v", cursor)
	}
	if got := m.Inventory().GetSlot(9); got.ItemID != 12 || got.Count != 1 {
		t.Fatalf("inactive window-0 content changed inventory to %#v", got)
	}

	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 5, StateID: 4, Slots: []slot.Slot{{Count: 1}, {Count: 2}}, CarriedItem: slot.Slot{Count: 3},
	})
	if got := m.Container().SlotCount(); got != 2 {
		t.Fatalf("slot count = %d, want 2", got)
	}
	if cursor := m.Cursor(); cursor == nil || cursor.Count != 3 {
		t.Fatalf("cursor = %#v, want count 3", cursor)
	}
	if state := c.player.StateID(); state != 4 {
		t.Fatalf("state ID = %d, want 4", state)
	}
}

func TestManagerTracksCursorSlotAndClearsClosedContainer(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	title := chat.Text("Test ").SetColor(chat.Gold).Append(chat.Text("Menu"))
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{
		WindowID: 5, WindowType: 12, WindowTitle: title,
	})
	if got := m.CurrentMenuTitle(); m.CurrentMenuType() != 12 ||
		got.Color != chat.Gold || got.ClearString() != "Test Menu" {
		t.Fatalf("menu metadata = (%d, %#v), want type 12 and complete title component", m.CurrentMenuType(), got)
	}
	c.handler.HandlePacket(context.Background(), &gameclient.ContainerSetSlot{
		ContainerID: -1, Slot: -1, StateID: 7, ItemStack: slot.Slot{Count: 6},
	})
	if cursor := m.Cursor(); cursor == nil || cursor.Count != 6 {
		t.Fatalf("cursor = %#v, want count 6", cursor)
	}
	if got := m.Container().SlotCount(); got != 0 {
		t.Fatalf("cursor sentinel became a regular slot; slot count = %d", got)
	}

	c.handler.HandlePacket(context.Background(), &gameclient.CloseContainer{WindowID: 5})
	if m.CurrentContainerID() != -1 || m.CurrentMenuType() != -1 || m.CurrentMenuTitle().ClearString() != "" ||
		m.Container() != nil || m.Cursor() != nil {
		t.Fatalf("closed manager retained state: id=%d type=%d title=%#v container=%#v cursor=%#v",
			m.CurrentContainerID(), m.CurrentMenuType(), m.CurrentMenuTitle(), m.Container(), m.Cursor())
	}
}

func TestManagerTracksDedicatedCursorAndPlayerInventoryPackets(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m

	c.handler.HandlePacket(context.Background(), &gameclient.SetCursorItem{
		CarriedItem: slot.Slot{ItemID: 7, Count: 3},
	})
	c.handler.HandlePacket(context.Background(), &gameclient.SetPlayerInventory{
		Slot: 0, Data: slot.Slot{ItemID: 8, Count: 2},
	})
	if cursor := m.Cursor(); cursor == nil || cursor.ItemID != 7 || cursor.Count != 3 {
		t.Fatalf("dedicated cursor packet produced %#v", cursor)
	}
	if got := m.Inventory().GetSlot(36); got.ItemID != 8 || got.Count != 2 {
		t.Fatalf("standalone hotbar slot = %#v", got)
	}
}

func TestManagerPreservesExternalMenuPlayerUpdatesAfterCursorLifecycleCloses(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m

	const playerStart = 27
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
	slots := make([]slot.Slot, playerStart+36)
	slots[playerStart] = slot.Slot{ItemID: 7, Count: 1}
	slots[playerStart+35] = slot.Slot{ItemID: 8, Count: 2}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 5,
		StateID:  4,
		Slots:    slots,
		// Reproduce a server that publishes the destination but omits the
		// corresponding cursor-clear update until the menu lifecycle ends.
		CarriedItem: slot.Slot{ItemID: 7, Count: 1},
	})
	c.handler.HandlePacket(context.Background(), &gameclient.ContainerSetSlot{
		ContainerID: 5,
		StateID:     5,
		Slot:        playerStart + 1,
		ItemStack:   slot.Slot{ItemID: 9, Count: 3},
	})

	m.Close()

	if cursor := m.Cursor(); cursor != nil {
		t.Fatalf("closed cursor lifecycle retained %#v", cursor)
	}
	for inventorySlot, want := range map[int]slot.Slot{
		9:  {ItemID: 7, Count: 1},
		10: {ItemID: 9, Count: 3},
		44: {ItemID: 8, Count: 2},
	} {
		if got := m.Inventory().GetSlot(inventorySlot); got.ItemID != want.ItemID || got.Count != want.Count {
			t.Fatalf("inventory slot %d after close = %#v, want %#v", inventorySlot, got, want)
		}
	}
}

func TestManagerMirrorsStandalonePlayerUpdatesIntoOpenMenu(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m

	const playerStart = 27
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 5,
		Slots:    make([]slot.Slot, playerStart+36),
	})
	c.handler.HandlePacket(context.Background(), &gameclient.SetPlayerInventory{
		Slot: 0,
		Data: slot.Slot{ItemID: 8, Count: 2},
	})
	c.handler.HandlePacket(context.Background(), &gameclient.SetPlayerInventory{
		Slot: 9,
		Data: slot.Slot{ItemID: 9, Count: 3},
	})
	c.handler.HandlePacket(context.Background(), &gameclient.ContainerSetSlot{
		ContainerID: 0,
		Slot:        36,
		ItemStack:   slot.Slot{ItemID: 10, Count: 4},
	})
	c.handler.HandlePacket(context.Background(), &gameclient.ContainerSetSlot{
		ContainerID: 0,
		Slot:        9,
		ItemStack:   slot.Slot{ItemID: 11, Count: 5},
	})

	if got := m.Container().GetSlot(playerStart + 27); got.ItemID != 10 || got.Count != 4 {
		t.Fatalf("open-menu hotbar slot = %#v, want item 10 count 4", got)
	}
	if got := m.Container().GetSlot(playerStart); got.ItemID != 9 || got.Count != 3 {
		t.Fatalf("open-menu main slot = %#v, want item 9 count 3", got)
	}
}

func TestManagerInventoryIncludesCraftingResultSlot(t *testing.T) {
	m := NewManager(newInventoryTestClient())
	if got := m.Inventory().SlotCount(); got != 46 {
		t.Fatalf("inventory slot count = %d, want 46", got)
	}
}

func TestContainerClickUsesItsOwnStateIDWithoutValidationData(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 5, StateID: 4, Slots: []slot.Slot{{Count: 1}},
	})
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 9, Slots: []slot.Slot{{Count: 1}},
	})

	if err := m.Container().Click(0, 0, 0); err != nil {
		t.Fatal(err)
	}
	click, ok := c.writes[0].(*server.ContainerClick)
	if !ok {
		t.Fatalf("packet = %T, want *ContainerClick", c.writes[0])
	}
	if click.StateID != 4 {
		t.Fatalf("container click state ID = %d, want 4", click.StateID)
	}
	if len(click.ChangedSlots) != 0 || click.CarriedSlot.HasItem {
		t.Fatalf("container click contains validation data: %#v", click)
	}
	if got := m.Container().GetSlot(0); got.Count != 1 {
		t.Fatalf("click predicted local slot state: %#v", got)
	}
}

func TestManagerClickDoesNotSendHashedValidationOrPredict(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	slots := make([]slot.Slot, 46)
	slots[9] = slot.Slot{ItemID: 7, Count: 2}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 12, Slots: slots,
	})

	if err := m.Click(0, 9, 0, 0); err != nil {
		t.Fatal(err)
	}
	click := c.writes[0].(*server.ContainerClick)
	if click.StateID != 12 || len(click.ChangedSlots) != 0 || click.CarriedSlot.HasItem {
		t.Fatalf("manager click contains validation data: %#v", click)
	}
	if got := m.Inventory().GetSlot(9); got.Count != 2 {
		t.Fatalf("manager click predicted local state: %#v", got)
	}
}

func TestManagerLoginAndRespawnClearOpenContainerWithoutClientClose(t *testing.T) {
	for _, lifecyclePacket := range []gameclient.ClientboundPacket{
		&gameclient.Login{},
		&gameclient.Respawn{},
	} {
		c := newInventoryTestClient()
		m := NewManager(c)
		c.inventory = m
		c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
		c.handler.HandlePacket(context.Background(), &gameclient.SetCursorItem{
			CarriedItem: slot.Slot{ItemID: 7, Count: 1},
		})
		c.handler.HandlePacket(context.Background(), lifecyclePacket)

		title := m.CurrentMenuTitle()
		if m.CurrentContainerID() != -1 || m.CurrentMenuType() != -1 ||
			title.ClearString() != "" || m.Container() != nil || m.Cursor() != nil {
			t.Fatalf("%T left stale inventory state: id=%d type=%d title=%#v container=%#v cursor=%#v",
				lifecyclePacket, m.CurrentContainerID(), m.CurrentMenuType(), m.CurrentMenuTitle(), m.Container(), m.Cursor())
		}
		if len(c.writes) != 0 {
			t.Fatalf("%T sent an unnecessary client close packet: %#v", lifecyclePacket, c.writes)
		}
	}
}

func TestManagerLoginAndRespawnClearCursorWithoutOpenContainer(t *testing.T) {
	for _, lifecyclePacket := range []gameclient.ClientboundPacket{
		&gameclient.Login{},
		&gameclient.Respawn{},
	} {
		c := newInventoryTestClient()
		m := NewManager(c)
		c.inventory = m
		c.handler.HandlePacket(context.Background(), &gameclient.SetCursorItem{
			CarriedItem: slot.Slot{ItemID: 7, Count: 1},
		})

		c.handler.HandlePacket(context.Background(), lifecyclePacket)

		if cursor := m.Cursor(); cursor != nil {
			t.Fatalf("%T left cursor without an open container: %#v", lifecyclePacket, cursor)
		}
	}
}

func TestContainerSlotsReturnsCopy(t *testing.T) {
	c := NewContainerWithSize(nil, 0, 1)
	c.SetSlot(0, slot.Slot{Count: 2, RemoveComponent: []int32{1}})
	slots := c.Slots()
	slots[0].Count = 9
	slots[0].RemoveComponent[0] = 9
	if got := c.GetSlot(0).Count; got != 2 {
		t.Fatalf("mutating Slots result changed container count to %d", got)
	}
	if got := c.GetSlot(0).RemoveComponent[0]; got != 1 {
		t.Fatalf("mutating nested slot data changed component to %d", got)
	}
}

type inventoryTestClient struct {
	handler   *inventoryPacketHandler
	player    *inventoryTestPlayer
	inventory bot.InventoryHandler
	writes    []server.ServerboundPacket
}

func newInventoryTestClient() *inventoryTestClient {
	return &inventoryTestClient{
		handler: &inventoryPacketHandler{handlers: make(map[packetid.ClientboundPacketID][]func(context.Context, gameclient.ClientboundPacket))},
		player:  &inventoryTestPlayer{},
	}
}

func (c *inventoryTestClient) Connect(context.Context, string, *bot.ConnectOptions) error { return nil }
func (c *inventoryTestClient) HandleGame(context.Context) error                           { return nil }
func (c *inventoryTestClient) Close(context.Context) error                                { return nil }
func (c *inventoryTestClient) IsConnected() bool                                          { return true }
func (c *inventoryTestClient) WritePacket(_ context.Context, packet server.ServerboundPacket) error {
	c.writes = append(c.writes, packet)
	return nil
}
func (c *inventoryTestClient) PacketHandler() bot.PacketHandler { return c.handler }
func (c *inventoryTestClient) EventHandler() bot.EventHandler   { return inventoryEventHandler{} }
func (c *inventoryTestClient) World() bot.World                 { return nil }
func (c *inventoryTestClient) Inventory() bot.InventoryHandler  { return c.inventory }
func (c *inventoryTestClient) Player() bot.Player               { return c.player }

type inventoryPacketHandler struct {
	handlers map[packetid.ClientboundPacketID][]func(context.Context, gameclient.ClientboundPacket)
}

func (h *inventoryPacketHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(context.Context, gameclient.ClientboundPacket)) {
	h.handlers[id] = append(h.handlers[id], handler)
}
func (*inventoryPacketHandler) AddRawPacketHandler(packetid.ClientboundPacketID, func(context.Context, packet.Packet)) {
}
func (*inventoryPacketHandler) AddGenericPacketHandler(func(context.Context, gameclient.ClientboundPacket)) {
}
func (h *inventoryPacketHandler) HandlePacket(ctx context.Context, packet gameclient.ClientboundPacket) {
	for _, handler := range h.handlers[packet.PacketID()] {
		handler(ctx, packet)
	}
}

type inventoryEventHandler struct{}

func (inventoryEventHandler) PublishEvent(string, any) error         { return nil }
func (inventoryEventHandler) SubscribeEvent(string, func(any) error) {}

type inventoryTestPlayer struct {
	stateID  int32
	sequence int32
}

func (p *inventoryTestPlayer) StateID() int32                   { return p.stateID }
func (p *inventoryTestPlayer) UpdateStateID(id int32)           { p.stateID = id }
func (p *inventoryTestPlayer) Sequence() int32                  { return p.sequence }
func (p *inventoryTestPlayer) UpdateSequence(id int32)          { p.sequence = id }
func (*inventoryTestPlayer) Entity() bot.Entity                 { return nil }
func (*inventoryTestPlayer) FlyTo(mgl64.Vec3) error             { return nil }
func (*inventoryTestPlayer) WalkTo(mgl64.Vec3) error            { return nil }
func (*inventoryTestPlayer) LookAt(mgl64.Vec3) error            { return nil }
func (*inventoryTestPlayer) UpdateLocation()                    {}
func (*inventoryTestPlayer) BreakBlock(protocol.Position) error { return nil }
func (*inventoryTestPlayer) PlaceBlock(protocol.Position) error { return nil }
func (*inventoryTestPlayer) PlaceBlockWithArgs(protocol.Position, int32, mgl64.Vec3) error {
	return nil
}
func (*inventoryTestPlayer) OpenContainer(protocol.Position, int32) (bot.Container, error) {
	return nil, nil
}
func (*inventoryTestPlayer) UseItem(int8) error                     { return nil }
func (*inventoryTestPlayer) OpenMenu(string) (bot.Container, error) { return nil, nil }
func (*inventoryTestPlayer) Command(string) error                   { return nil }
func (*inventoryTestPlayer) Chat(string) error                      { return nil }
func (*inventoryTestPlayer) CheckServer()                           {}
