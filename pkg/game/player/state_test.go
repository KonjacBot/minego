package player

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/go-mc/chat"
	"github.com/KonjacBot/go-mc/data/entity"
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level/block"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
	gameclient "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestPlayerPositionAcknowledgesAbsoluteState(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p
	p.entity.SetPosition(mgl64.Vec3{10, 20, 30})
	p.entity.SetRotation(mgl64.Vec2{40, 50})

	c.handler.HandlePacket(context.Background(), &gameclient.PlayerPosition{
		ID: 7, X: 2, Y: 5, Z: 6, XRot: 3, YRot: 4, Flags: 0x01 | 0x08,
	})

	if len(c.writes) != 2 {
		t.Fatalf("packet writes = %d, want 2", len(c.writes))
	}
	move, ok := c.writes[1].(*server.MovePlayerPosRot)
	if !ok {
		t.Fatalf("second packet = %T, want *MovePlayerPosRot", c.writes[1])
	}
	if move.X != 12 || move.FeetY != 5 || move.Z != 6 || move.XRot != 43 || move.YRot != 4 {
		t.Fatalf("absolute acknowledgement = %#v", move)
	}
}

func TestInteractionsUseIncreasingSequences(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p

	if err := p.PlaceBlock(protocol.Position{}); err != nil {
		t.Fatal(err)
	}
	if err := p.UseItem(0); err != nil {
		t.Fatal(err)
	}
	first := c.writes[0].(*server.UseItemOn).Sequence
	second := c.writes[1].(*server.UseItem).Sequence
	if first <= 0 || second <= first {
		t.Fatalf("interaction sequences = %d, %d, want increasing", first, second)
	}
}

func TestBreakBlockCompletionWaitsForServerAcknowledgement(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p

	if err := p.BreakBlock(protocol.Position{1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	if len(c.writes) != 2 {
		t.Fatalf("break packet writes = %d, want start and finish", len(c.writes))
	}
	finish, ok := c.writes[1].(*server.PlayerAction)
	if !ok || finish.Status != 2 {
		t.Fatalf("finish packet = %#v, want status 2 PlayerAction", c.writes[1])
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	acknowledged := make(chan error, 1)
	go func() {
		acknowledged <- p.WaitForSequence(ctx, finish.Sequence)
	}()
	c.handler.HandlePacket(context.Background(), &gameclient.BlockChangedAck{Sequence: finish.Sequence})

	if err := <-acknowledged; err != nil {
		t.Fatalf("wait for finish acknowledgement: %v", err)
	}
	if got := p.AcknowledgedSequence(); got != finish.Sequence {
		t.Fatalf("acknowledged sequence = %d, want %d", got, finish.Sequence)
	}
}

func TestBlockActionMethodsReturnTheExactPacketSequence(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p
	position := protocol.Position{1, 2, 3}

	start, err := p.StartBreakingBlock(context.Background(), position, 4)
	if err != nil {
		t.Fatal(err)
	}
	finish, err := p.FinishBreakingBlock(context.Background(), position, 4)
	if err != nil {
		t.Fatal(err)
	}
	placed, err := p.PlaceBlockWithArgsContext(context.Background(), position, 5, mgl64.Vec3{0.25, 0.5, 0.75})
	if err != nil {
		t.Fatal(err)
	}

	if start <= 0 || finish <= start || placed <= finish {
		t.Fatalf("returned sequences = start %d, finish %d, place %d; want strictly increasing", start, finish, placed)
	}
	startPacket := c.writes[0].(*server.PlayerAction)
	finishPacket := c.writes[1].(*server.PlayerAction)
	placePacket := c.writes[2].(*server.UseItemOn)
	if startPacket.Status != 0 || startPacket.Face != 4 || startPacket.Sequence != start {
		t.Fatalf("start packet = %#v", startPacket)
	}
	if finishPacket.Status != 2 || finishPacket.Face != 4 || finishPacket.Sequence != finish {
		t.Fatalf("finish packet = %#v", finishPacket)
	}
	if placePacket.Face != 5 || placePacket.Sequence != placed {
		t.Fatalf("place packet = %#v", placePacket)
	}
}

func TestFlyToUsesAirborneMovementState(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p
	p.abilities = 0x04

	start := mgl64.Vec3{-474.5, 185, 5625.5}
	target := mgl64.Vec3{-478, 184, 5622}
	p.entity.SetPosition(start)
	var groundedMovement bool
	c.writeHook = func(packet server.ServerboundPacket) {
		move, ok := packet.(*server.MovePlayerPos)
		if !ok {
			return
		}
		if move.Flags&0x01 != 0 {
			groundedMovement = true
		}
	}

	if err := p.FlyTo(target); err != nil {
		t.Fatal(err)
	}
	if groundedMovement {
		t.Fatal("FlyTo marked airborne movement as on-ground")
	}
}

func TestFlyToPreservesWorkingFiveBlockWarehouseCadence(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p
	p.abilities = 0x04

	start := mgl64.Vec3{-474.5, 185, 5625.5}
	target := mgl64.Vec3{-478, 184, 5622}
	p.entity.SetPosition(start)

	if err := p.FlyTo(target); err != nil {
		t.Fatal(err)
	}

	var moves []mgl64.Vec3
	for _, packet := range c.writes {
		if move, ok := packet.(*server.MovePlayerPos); ok {
			moves = append(moves, mgl64.Vec3{move.X, move.FeetY, move.Z})
		}
	}
	if len(moves) != 1 {
		t.Fatalf("warehouse flight packets = %d; want 1 using the working five-block cadence", len(moves))
	}
	if got := moves[0].Sub(start).Len(); math.Abs(got-5) > 0.001 {
		t.Fatalf("first warehouse flight segment = %.3f blocks; want 5", got)
	}
}

func TestFlyToRejectsImmediateServerCorrectionAtLargeCoordinates(t *testing.T) {
	c := newStateTestClient()
	p := New(c)
	c.player = p
	p.abilities = 0x04

	start := mgl64.Vec3{-474.5, 185, 5625.5}
	target := mgl64.Vec3{-478, 184, 5622}
	p.entity.SetPosition(start)
	var correction sync.Once
	c.writeHook = func(packet server.ServerboundPacket) {
		if _, ok := packet.(*server.MovePlayerPos); ok {
			correction.Do(func() { p.entity.SetPosition(start) })
		}
	}

	if err := p.FlyTo(target); err == nil {
		t.Fatal("FlyTo accepted a server correction because large coordinates hid the movement delta")
	}
}

func TestWaitForContainerReturnsNewInitializedWindow(t *testing.T) {
	c := newStateTestClient()
	inventory := &waitInventory{id: -1}
	c.inventory = inventory
	p := New(c)
	c.player = p

	go func() {
		time.Sleep(20 * time.Millisecond)
		inventory.mu.Lock()
		inventory.id = 5
		inventory.container = openContainerTestContainer{slots: make([]slot.Slot, 1)}
		inventory.mu.Unlock()
	}()
	container, err := p.waitForContainer(-1)
	if err != nil {
		t.Fatal(err)
	}
	if container == nil || container.SlotCount() != 1 {
		t.Fatalf("container = %#v, want initialized window", container)
	}
}

func TestIsWalkableFailsClosedForUnknownBlocks(t *testing.T) {
	w := pathTestWorld{err: errors.New("chunk not loaded")}
	if isWalkable(w, protocol.Position{}) {
		t.Fatal("isWalkable returned true for an unloaded chunk")
	}
}

func TestIsWalkableRequiresSupport(t *testing.T) {
	pos := protocol.Position{1, 2, 3}
	w := pathTestWorld{blocks: map[protocol.Position]block.Block{
		pos:                          block.Air{},
		{pos[0], pos[1] + 1, pos[2]}: block.Air{},
		{pos[0], pos[1] - 1, pos[2]}: block.Stone{},
	}}
	if !isWalkable(w, pos) {
		t.Fatal("isWalkable returned false for supported air blocks")
	}
	w.blocks[protocol.Position{pos[0], pos[1] - 1, pos[2]}] = block.Air{}
	if isWalkable(w, pos) {
		t.Fatal("isWalkable returned true without a supporting block")
	}
}

type stateTestClient struct {
	handler   *statePacketHandler
	player    bot.Player
	inventory bot.InventoryHandler
	writes    []server.ServerboundPacket
	writeHook func(server.ServerboundPacket)
}

func newStateTestClient() *stateTestClient {
	return &stateTestClient{handler: &statePacketHandler{handlers: make(map[packetid.ClientboundPacketID][]func(context.Context, gameclient.ClientboundPacket))}}
}

func (c *stateTestClient) Connect(context.Context, string, *bot.ConnectOptions) error { return nil }
func (c *stateTestClient) HandleGame(context.Context) error                           { return nil }
func (c *stateTestClient) Close(context.Context) error                                { return nil }
func (c *stateTestClient) IsConnected() bool                                          { return true }
func (c *stateTestClient) WritePacket(_ context.Context, packet server.ServerboundPacket) error {
	c.writes = append(c.writes, packet)
	if c.writeHook != nil {
		c.writeHook(packet)
	}
	return nil
}
func (c *stateTestClient) PacketHandler() bot.PacketHandler { return c.handler }
func (c *stateTestClient) EventHandler() bot.EventHandler   { return stateEventHandler{} }
func (c *stateTestClient) World() bot.World                 { return nil }
func (c *stateTestClient) Inventory() bot.InventoryHandler  { return c.inventory }
func (c *stateTestClient) Player() bot.Player               { return c.player }

type statePacketHandler struct {
	handlers map[packetid.ClientboundPacketID][]func(context.Context, gameclient.ClientboundPacket)
	generic  []func(context.Context, gameclient.ClientboundPacket)
}

func (h *statePacketHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(context.Context, gameclient.ClientboundPacket)) {
	h.handlers[id] = append(h.handlers[id], handler)
}
func (h *statePacketHandler) AddRawPacketHandler(packetid.ClientboundPacketID, func(context.Context, pk.Packet)) {
}
func (h *statePacketHandler) AddGenericPacketHandler(handler func(context.Context, gameclient.ClientboundPacket)) {
	h.generic = append(h.generic, handler)
}
func (h *statePacketHandler) HandlePacket(ctx context.Context, packet gameclient.ClientboundPacket) {
	for _, handler := range h.generic {
		handler(ctx, packet)
	}
	for _, handler := range h.handlers[packet.PacketID()] {
		handler(ctx, packet)
	}
}

type stateEventHandler struct{}

func (stateEventHandler) PublishEvent(string, any) error         { return nil }
func (stateEventHandler) SubscribeEvent(string, func(any) error) {}

type waitInventory struct {
	mu        sync.RWMutex
	id        int32
	container bot.Container
}

func (i *waitInventory) Inventory() bot.Container { return nil }
func (i *waitInventory) Container() bot.Container {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.container
}
func (i *waitInventory) CurrentContainerID() int32 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.id
}
func (*waitInventory) CurrentMenuType() int32                 { return -1 }
func (*waitInventory) CurrentMenuTitle() chat.Message         { return chat.Message{} }
func (*waitInventory) Cursor() *slot.Slot                     { return nil }
func (*waitInventory) Click(int32, int16, int32, int32) error { return nil }
func (*waitInventory) Close()                                 {}

type pathTestWorld struct {
	blocks map[protocol.Position]block.Block
	err    error
}

func (w pathTestWorld) GetBlock(pos protocol.Position) (block.Block, error) {
	if w.err != nil {
		return nil, w.err
	}
	value, ok := w.blocks[pos]
	if !ok {
		return nil, errors.New("block not loaded")
	}
	return value, nil
}
func (pathTestWorld) SetBlock(protocol.Position, block.Block) error { return nil }
func (pathTestWorld) GetNearbyBlocks(protocol.Position, int32) ([]block.Block, error) {
	return nil, nil
}
func (pathTestWorld) FindNearbyBlock(protocol.Position, int32, block.Block) (protocol.Position, error) {
	return protocol.Position{}, errors.New("not implemented")
}
func (pathTestWorld) Entities() []bot.Entity               { return nil }
func (pathTestWorld) GetEntity(int32) bot.Entity           { return nil }
func (pathTestWorld) GetNearbyEntities(int32) []bot.Entity { return nil }
func (pathTestWorld) GetEntitiesByType(entity.ID) []bot.Entity {
	return nil
}
