package inventory

import (
	"context"
	"testing"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

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

	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 4, StateID: 3, Slots: []slot.Slot{{Count: 1}}, CarriedItem: slot.Slot{Count: 9},
	})
	if got := m.Container().SlotCount(); got != 0 {
		t.Fatalf("stale window changed slot count to %d", got)
	}
	if cursor := m.Cursor(); cursor != nil {
		t.Fatalf("stale window changed cursor to %#v", cursor)
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
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5})
	c.handler.HandlePacket(context.Background(), &gameclient.ContainerSetSlot{
		ContainerID: -1, Slot: -1, StateID: 7, ItemStack: slot.Slot{Count: 6},
	})
	if cursor := m.Cursor(); cursor == nil || cursor.Count != 6 {
		t.Fatalf("cursor = %#v, want count 6", cursor)
	}

	c.handler.HandlePacket(context.Background(), &gameclient.CloseContainer{WindowID: 5})
	if m.CurrentContainerID() != -1 || m.Container() != nil || m.Cursor() != nil {
		t.Fatalf("closed manager retained state: id=%d container=%#v cursor=%#v", m.CurrentContainerID(), m.Container(), m.Cursor())
	}
}

func TestContainerClickUsesItsOwnStateID(t *testing.T) {
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

func TestManagerClickSendsPredictionAndUpdatesLocalCache(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 12, Slots: make([]slot.Slot, 46),
		CarriedItem: slot.Slot{ItemID: 7, Count: 3},
	})

	if err := m.Click(0, 9, 0, 0); err != nil {
		t.Fatal(err)
	}
	click := c.writes[0].(*server.ContainerClick)
	if click.StateID != 12 || len(click.ChangedSlots) != 1 || click.ChangedSlots[0].Slot != 9 {
		t.Fatalf("click prediction = %#v", click)
	}
	if !click.ChangedSlots[0].SlotData.HasItem || click.ChangedSlots[0].SlotData.ItemCount != 3 || click.CarriedSlot.HasItem {
		t.Fatalf("click hashed slots = changed %#v, carried %#v", click.ChangedSlots[0].SlotData, click.CarriedSlot)
	}
	if got := m.Inventory().GetSlot(9); got.ItemID != 7 || got.Count != 3 {
		t.Fatalf("predicted local slot = %#v", got)
	}
	if cursor := m.Cursor(); cursor == nil || cursor.Count != 0 {
		t.Fatalf("predicted cursor = %#v", cursor)
	}
}

func TestManagerShiftClickUsesGenericContainerLayout(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5, WindowType: 2})
	slots := make([]slot.Slot, 63)
	slots[0] = slot.Slot{ItemID: 7, Count: 4}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{WindowID: 5, StateID: 3, Slots: slots})

	if err := m.Click(5, 0, 1, 0); err != nil {
		t.Fatal(err)
	}
	click := c.writes[0].(*server.ContainerClick)
	if len(click.ChangedSlots) != 2 || click.ChangedSlots[0].Slot != 0 || click.ChangedSlots[1].Slot != 62 {
		t.Fatalf("shift-click changed slots = %#v", click.ChangedSlots)
	}
	if got := m.Container().GetSlot(0); got.Count != 0 {
		t.Fatalf("source slot still contains %#v", got)
	}
	if got := m.Container().GetSlot(62); got.ItemID != 7 || got.Count != 4 {
		t.Fatalf("predicted destination = %#v", got)
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
func (*inventoryPacketHandler) AddRawPacketHandler(packetid.ClientboundPacketID, func(context.Context, pk.Packet)) {
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
