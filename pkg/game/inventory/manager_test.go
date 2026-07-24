package inventory

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

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
	if snapshot := m.Snapshot(); snapshot.Revision != 1 || snapshot.AuthoritativeRevision != 1 {
		t.Fatalf("stale content advanced revisions to (%d, %d)", snapshot.Revision, snapshot.AuthoritativeRevision)
	}
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
	if got := m.Container().SlotCount(); got != 0 {
		t.Fatalf("cursor sentinel became a regular slot; slot count = %d", got)
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

func TestManagerSnapshotIsDetachedAndAtomic(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	slots := make([]slot.Slot, 46)
	slots[9] = slot.Slot{ItemID: 7, Count: 2, RemoveComponent: []int32{1}}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 3, Slots: slots, CarriedItem: slot.Slot{ItemID: 8, Count: 1},
	})

	snapshot := m.Snapshot()
	if snapshot.WindowID != 0 || snapshot.StateID != 3 || !snapshot.Ready || snapshot.PlayerSlotStart != 9 {
		t.Fatalf("snapshot metadata = %#v", snapshot)
	}
	snapshot.ContainerSlots[9].Count = 9
	snapshot.ContainerSlots[9].RemoveComponent[0] = 9
	snapshot.PlayerInventorySlots[9].Count = 8
	snapshot.Cursor.Count = 7

	again := m.Snapshot()
	if again.ContainerSlots[9].Count != 2 || again.ContainerSlots[9].RemoveComponent[0] != 1 {
		t.Fatalf("mutating container snapshot changed manager state: %#v", again.ContainerSlots[9])
	}
	if again.PlayerInventorySlots[9].Count != 2 || again.Cursor.Count != 1 {
		t.Fatalf("mutating inventory snapshot changed manager state: inventory=%#v cursor=%#v", again.PlayerInventorySlots[9], again.Cursor)
	}
}

func TestManagerLocalAndAuthoritativeRevisions(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	slots := make([]slot.Slot, 46)
	slots[9] = slot.Slot{ItemID: 7, Count: 1}
	content := &gameclient.SetContainerContent{WindowID: 0, StateID: 4, Slots: slots}
	c.handler.HandlePacket(context.Background(), content)
	server := m.Snapshot()
	if server.Revision != 1 || server.AuthoritativeRevision != 1 {
		t.Fatalf("server revisions = (%d, %d), want (1, 1)", server.Revision, server.AuthoritativeRevision)
	}

	if err := m.ClickContext(context.Background(), 0, 9, 0, 0); err != nil {
		t.Fatal(err)
	}
	predicted := m.Snapshot()
	if predicted.Revision != 2 || predicted.AuthoritativeRevision != 1 || predicted.Cursor.Count != 1 {
		t.Fatalf("prediction snapshot = %#v", predicted)
	}

	// An equal authoritative packet is still a handled protocol observation.
	c.handler.HandlePacket(context.Background(), content)
	corrected := m.Snapshot()
	if corrected.Revision != 3 || corrected.AuthoritativeRevision != 2 {
		t.Fatalf("correction revisions = (%d, %d), want (3, 2)", corrected.Revision, corrected.AuthoritativeRevision)
	}
	if corrected.ContainerSlots[9].Count != 1 || corrected.Cursor.Count != 0 {
		t.Fatalf("server correction did not overwrite prediction: %#v", corrected)
	}
}

func TestManagerClickTransactionDistinguishesNoOpAndUnsupportedPrediction(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 4, Slots: make([]slot.Slot, 46),
	})
	baseline := m.Snapshot()

	noOp, err := m.ClickTransaction(context.Background(), 0, 9, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !noOp.Complete || noOp.Changed {
		t.Fatalf("complete no-op result = %#v", noOp)
	}
	if noOp.Revision != baseline.Revision || noOp.AuthoritativeRevision != baseline.AuthoritativeRevision {
		t.Fatalf("complete no-op revisions = %#v, baseline = %#v", noOp, baseline)
	}

	unsupported, err := m.ClickTransaction(context.Background(), 0, 9, 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if unsupported.Complete || unsupported.Changed {
		t.Fatalf("unsupported prediction result = %#v", unsupported)
	}
	if unsupported.Revision != baseline.Revision || unsupported.AuthoritativeRevision != baseline.AuthoritativeRevision {
		t.Fatalf("unsupported revisions = %#v, baseline = %#v", unsupported, baseline)
	}
}

func TestManagerMenuEpochChangesWhenWindowIDIsReused(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5, WindowType: 2})
	first := m.Snapshot()
	c.handler.HandlePacket(context.Background(), &gameclient.CloseContainer{WindowID: 5})
	c.handler.HandlePacket(context.Background(), &gameclient.OpenScreen{WindowID: 5, WindowType: 2})
	second := m.Snapshot()

	if first.WindowID != second.WindowID || first.MenuEpoch == 0 || second.MenuEpoch <= first.MenuEpoch {
		t.Fatalf("reused window epochs = first %#v, second %#v", first, second)
	}
}

func TestManagerWaitDoesNotLoseUpdate(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	after := m.Snapshot().Revision

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result := make(chan Snapshot, 1)
	errResult := make(chan error, 1)
	go func() {
		snapshot, err := m.Wait(ctx, after, func(snapshot Snapshot) bool { return snapshot.Ready })
		result <- snapshot
		errResult <- err
	}()
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 6, Slots: make([]slot.Slot, 46),
	})

	if err := <-errResult; err != nil {
		t.Fatal(err)
	}
	if snapshot := <-result; snapshot.Revision <= after || !snapshot.Ready || snapshot.StateID != 6 {
		t.Fatalf("wait snapshot = %#v", snapshot)
	}
}

func TestManagerWaitPredicateMayReadManagerSnapshot(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 6, Slots: make([]slot.Slot, 46),
	})

	predicateEntered := make(chan struct{})
	writerAttempted := make(chan struct{})
	writerDone := make(chan struct{})
	go func() {
		<-predicateEntered
		close(writerAttempted)
		c.handler.HandlePacket(context.Background(), &gameclient.SetCursorItem{})
		close(writerDone)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	type waitResult struct {
		snapshot Snapshot
		err      error
	}
	result := make(chan waitResult, 1)
	go func() {
		snapshot, err := m.Wait(ctx, 0, func(Snapshot) bool {
			close(predicateEntered)
			<-writerAttempted
			// Give the writer time to queue before recursively taking an RLock.
			time.Sleep(20 * time.Millisecond)
			_ = m.Snapshot()
			return true
		})
		result <- waitResult{snapshot: snapshot, err: err}
	}()

	select {
	case got := <-result:
		if got.err != nil || got.snapshot.StateID != 6 {
			t.Fatalf("Wait result = (%#v, %v)", got.snapshot, got.err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Wait predicate deadlocked while reading Manager snapshot")
	}
	select {
	case <-writerDone:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("writer remained blocked by Wait predicate")
	}
}

func TestManagerClickContextCanceledBeforeWrite(t *testing.T) {
	c := newInventoryTestClient()
	m := NewManager(c)
	c.inventory = m
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{
		WindowID: 0, StateID: 2, Slots: make([]slot.Slot, 46),
		CarriedItem: slot.Slot{ItemID: 7, Count: 1},
	})
	before := m.Snapshot()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := m.ClickContext(ctx, 0, 9, 0, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ClickContext error = %v, want context.Canceled", err)
	}
	if len(c.writtenPackets()) != 0 {
		t.Fatal("canceled click started a network side effect")
	}
	if after := m.Snapshot(); after.Revision != before.Revision || after.ContainerSlots[9].Count != 0 {
		t.Fatalf("canceled click changed state: before=%#v after=%#v", before, after)
	}
}

func TestManagerSerializesConcurrentClickTransactions(t *testing.T) {
	c := newInventoryTestClient()
	entered := make(chan struct{}, 2)
	releaseFirst := make(chan struct{})
	var hookMu sync.Mutex
	writes := 0
	c.writeHook = func(ctx context.Context, _ server.ServerboundPacket) error {
		hookMu.Lock()
		writes++
		write := writes
		hookMu.Unlock()
		entered <- struct{}{}
		if write == 1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-releaseFirst:
			}
		}
		return nil
	}
	m := NewManager(c)
	c.inventory = m
	slots := make([]slot.Slot, 46)
	slots[9] = slot.Slot{ItemID: 7, Count: 1}
	c.handler.HandlePacket(context.Background(), &gameclient.SetContainerContent{WindowID: 0, StateID: 5, Slots: slots})

	errs := make(chan error, 2)
	go func() { errs <- m.ClickContext(context.Background(), 0, 9, 0, 0) }()
	<-entered
	go func() { errs <- m.ClickContext(context.Background(), 0, 10, 0, 0) }()
	select {
	case <-entered:
		t.Fatal("second click began its write before the first transaction committed")
	case <-time.After(50 * time.Millisecond):
	}
	close(releaseFirst)
	for range 2 {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	packets := c.writtenPackets()
	if len(packets) != 2 {
		t.Fatalf("writes = %d, want 2", len(packets))
	}
	second := packets[1].(*server.ContainerClick)
	if len(second.ChangedSlots) != 1 || second.ChangedSlots[0].Slot != 10 || second.CarriedSlot.HasItem {
		t.Fatalf("second click did not use first prediction: %#v", second)
	}
	if snapshot := m.Snapshot(); snapshot.ContainerSlots[9].Count != 0 || snapshot.ContainerSlots[10].Count != 1 || snapshot.Cursor.Count != 0 {
		t.Fatalf("serialized prediction state = %#v", snapshot)
	}
}

type inventoryTestClient struct {
	mu        sync.Mutex
	handler   *inventoryPacketHandler
	player    *inventoryTestPlayer
	inventory bot.InventoryHandler
	writes    []server.ServerboundPacket
	writeHook func(context.Context, server.ServerboundPacket) error
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
func (c *inventoryTestClient) WritePacket(ctx context.Context, packet server.ServerboundPacket) error {
	if c.writeHook != nil {
		if err := c.writeHook(ctx, packet); err != nil {
			return err
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writes = append(c.writes, packet)
	return nil
}
func (c *inventoryTestClient) writtenPackets() []server.ServerboundPacket {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]server.ServerboundPacket(nil), c.writes...)
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
