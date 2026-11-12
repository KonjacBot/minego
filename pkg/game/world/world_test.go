package world

import (
	"bytes"
	"context"
	"testing"

	"github.com/go-gl/mathgl/mgl64"

	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/bot"
	"github.com/KonjacBot/minego/v2/pkg/bus"
	gameclient "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
)

func TestEntityCombinedMoveUpdatesPositionAndRotation(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w

	c.handler.DispatchPacket(context.Background(), &gameclient.AddEntity{
		ID: 1, X: 1, Y: 2, Z: 3, Yaw: pk.Angle(64), Pitch: pk.Angle(32),
	})
	c.handler.DispatchPacket(context.Background(), &gameclient.UpdateEntityPositionAndRotation{
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

// 原版只在換維度時清掉實體;同一個維度的重生留著,伺服器不會重送還在追蹤的展示框。
func TestRespawnClearsEntitiesOnlyOnDimensionChange(t *testing.T) {
	spawn := func(dim string) gameclient.CommonPlayerSpawnInfo {
		return gameclient.CommonPlayerSpawnInfo{Dimension: dim}
	}
	tests := []struct {
		name    string
		respawn string
		kept    bool
	}{
		{"same dimension", "minecraft:overworld", true},
		{"other dimension", "minecraft:the_nether", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newWorldTestClient()
			w := NewWorld(c)
			c.world = w
			c.handler.DispatchPacket(context.Background(), &gameclient.Login{CommonPlayerSpawnInfo: spawn("minecraft:overworld")})
			c.handler.DispatchPacket(context.Background(), &gameclient.AddEntity{ID: 1})
			c.handler.DispatchPacket(context.Background(), &gameclient.Respawn{CommonPlayerSpawnInfo: spawn(tt.respawn)})
			if kept := w.GetEntity(1) != nil; kept != tt.kept {
				t.Fatalf("entity kept = %v, want %v", kept, tt.kept)
			}
		})
	}
}

func TestLoginClearsEntities(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	c.handler.DispatchPacket(context.Background(), &gameclient.AddEntity{ID: 1})
	c.handler.DispatchPacket(context.Background(), &gameclient.Login{})
	if entity := w.GetEntity(1); entity != nil {
		t.Fatalf("entity survived login: %#v", entity)
	}
}

func TestAbsoluteEntityAndVehiclePositionUpdates(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	c.handler.DispatchPacket(context.Background(), &gameclient.AddEntity{ID: 1})
	c.handler.DispatchPacket(context.Background(), &gameclient.TeleportEntity{
		EntityID: 1,
		X:        10,
		Y:        64,
		Z:        -4,
		Yaw:      90,
		Pitch:    15,
	})
	got := w.GetEntity(1)
	if got == nil || got.Position() != (mgl64.Vec3{10, 64, -4}) ||
		got.Rotation() != (mgl64.Vec2{90, 15}) {
		t.Fatalf("teleported entity = %#v", got)
	}

	c.handler.DispatchPacket(context.Background(), &gameclient.SynchronizeVehiclePosition{
		EntityID: 1,
		X:        20,
		Y:        65,
		Z:        3,
		Yaw:      180,
		Pitch:    -5,
	})
	if got.Position() != (mgl64.Vec3{20, 65, 3}) ||
		got.Rotation() != (mgl64.Vec2{180, -5}) {
		t.Fatalf("synchronized vehicle = position %v rotation %v", got.Position(), got.Rotation())
	}
}

type worldTestClient struct {
	handler *bus.Bus
	events  *worldEventHandler
	world   bot.World
}

func newWorldTestClient() *worldTestClient {
	c := &worldTestClient{handler: bus.New(), events: &worldEventHandler{}}
	c.handler.On(func(context.Context, EntityAddEvent) { c.events.added++ })
	return c
}

func (c *worldTestClient) Connect(context.Context, string, *bot.ConnectOptions) error { return nil }
func (c *worldTestClient) HandleGame(context.Context) error                           { return nil }
func (c *worldTestClient) Close(context.Context) error                                { return nil }
func (c *worldTestClient) IsConnected() bool                                          { return true }
func (c *worldTestClient) WritePacket(context.Context, server.ServerboundPacket) error {
	return nil
}
func (c *worldTestClient) Bus() *bus.Bus                   { return c.handler }
func (c *worldTestClient) World() bot.World                { return c.world }
func (c *worldTestClient) Inventory() bot.InventoryHandler { return nil }
func (c *worldTestClient) Player() bot.Player              { return nil }

type worldEventHandler struct {
	added int
}

func TestAddEntityKeepsData(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	// Item frames send their facing direction as the AddEntity data (2 = north).
	c.handler.DispatchPacket(context.Background(), &gameclient.AddEntity{ID: 1, Data: 2})
	entity := w.GetEntity(1)
	if entity == nil {
		t.Fatal("entity was not added")
	}
	if got := entity.Data(); got != 2 {
		t.Fatalf("data = %d, want 2", got)
	}
}

func TestSetEquipmentAppliesToEntity(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	c.handler.DispatchPacket(context.Background(), &gameclient.AddEntity{ID: 300})

	// EntityID 300;主手 item 5、頭盔 item 1112,各 1 個(實機線路格式)。
	wire := []byte{0xac, 0x02, 0x80, 0x01, 0x05, 0x00, 0x00, 0x05, 0x01, 0xd8, 0x08, 0x00, 0x00}
	var p gameclient.SetEquipment
	if _, err := p.ReadFrom(bytes.NewReader(wire)); err != nil {
		t.Fatalf("ReadFrom(): %v", err)
	}
	c.handler.DispatchPacket(context.Background(), &p)

	entity := w.GetEntity(300)
	if entity == nil {
		t.Fatal("entity was not added")
	}
	equipment := entity.Equipment()
	if len(equipment) != 2 {
		t.Fatalf("Equipment() = %#v, want 2 slots", equipment)
	}
	if got := equipment[0]; got.ItemID != 5 || got.Count != 1 {
		t.Fatalf("main hand = %#v, want item 5 x1", got)
	}
	if got := equipment[5]; got.ItemID != 1112 || got.Count != 1 {
		t.Fatalf("head = %#v, want item 1112 x1", got)
	}
}
