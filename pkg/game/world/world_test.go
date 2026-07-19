package world

import (
	"context"
	"testing"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/bot"
	gameclient "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
)

func TestEntityCombinedMoveUpdatesPositionAndRotation(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w

	c.handler.HandlePacket(context.Background(), &gameclient.AddEntity{
		ID: 1, X: 1, Y: 2, Z: 3, Yaw: pk.Angle(64), Pitch: pk.Angle(32),
	})
	c.handler.HandlePacket(context.Background(), &gameclient.UpdateEntityPositionAndRotation{
		EntityID: 1, DeltaX: 4096, DeltaY: -4096, DeltaZ: 2048, Yaw: pk.Angle(-128), Pitch: pk.Angle(64),
	})

	entity := w.GetEntity(1)
	if entity == nil {
		t.Fatal("entity was not added")
	}
	position := entity.Position()
	rotation := entity.Rotation()
	if position[0] != 2 || position[1] != 1 || position[2] != 3.5 {
		t.Fatalf("position = %v, want [2 1 3.5]", position)
	}
	if rotation[0] != -180 || rotation[1] != 90 {
		t.Fatalf("rotation = %v, want [-180 90]", rotation)
	}
	if c.events.added != 1 {
		t.Fatalf("entity add events = %d, want 1", c.events.added)
	}
}

func TestRespawnClearsEntities(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	c.handler.HandlePacket(context.Background(), &gameclient.AddEntity{ID: 1})
	c.handler.HandlePacket(context.Background(), &gameclient.Respawn{})
	if entity := w.GetEntity(1); entity != nil {
		t.Fatalf("entity survived respawn: %#v", entity)
	}
}

type worldTestClient struct {
	handler *worldPacketHandler
	events  *worldEventHandler
	world   bot.World
}

func newWorldTestClient() *worldTestClient {
	return &worldTestClient{
		handler: &worldPacketHandler{handlers: make(map[packetid.ClientboundPacketID][]func(context.Context, gameclient.ClientboundPacket))},
		events:  &worldEventHandler{},
	}
}

func (c *worldTestClient) Connect(context.Context, string, *bot.ConnectOptions) error { return nil }
func (c *worldTestClient) HandleGame(context.Context) error                           { return nil }
func (c *worldTestClient) Close(context.Context) error                                { return nil }
func (c *worldTestClient) IsConnected() bool                                          { return true }
func (c *worldTestClient) WritePacket(context.Context, server.ServerboundPacket) error {
	return nil
}
func (c *worldTestClient) PacketHandler() bot.PacketHandler { return c.handler }
func (c *worldTestClient) EventHandler() bot.EventHandler   { return c.events }
func (c *worldTestClient) World() bot.World                 { return c.world }
func (c *worldTestClient) Inventory() bot.InventoryHandler  { return nil }
func (c *worldTestClient) Player() bot.Player               { return nil }

type worldPacketHandler struct {
	handlers map[packetid.ClientboundPacketID][]func(context.Context, gameclient.ClientboundPacket)
}

func (h *worldPacketHandler) AddPacketHandler(id packetid.ClientboundPacketID, handler func(context.Context, gameclient.ClientboundPacket)) {
	h.handlers[id] = append(h.handlers[id], handler)
}
func (*worldPacketHandler) AddRawPacketHandler(packetid.ClientboundPacketID, func(context.Context, pk.Packet)) {
}
func (*worldPacketHandler) AddGenericPacketHandler(func(context.Context, gameclient.ClientboundPacket)) {
}
func (h *worldPacketHandler) HandlePacket(ctx context.Context, packet gameclient.ClientboundPacket) {
	for _, handler := range h.handlers[packet.PacketID()] {
		handler(ctx, packet)
	}
}

type worldEventHandler struct {
	added int
}

func (h *worldEventHandler) PublishEvent(event string, _ any) error {
	if event == (EntityAddEvent{}).EventID() {
		h.added++
	}
	return nil
}
func (*worldEventHandler) SubscribeEvent(string, func(any) error) {}
