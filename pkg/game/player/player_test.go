package player

import (
	"context"
	"errors"
	"testing"

	"github.com/KonjacBot/go-mc/data/entity"
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level/block"
	"github.com/KonjacBot/go-mc/level/item"
	"github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
	gameclient "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func TestOpenContainerRefusesAirBlock(t *testing.T) {
	c := &openContainerTestClient{
		world:     openContainerTestWorld{blk: block.Air{}},
		inventory: openContainerTestInventory{id: 1, container: openContainerTestContainer{slots: make([]slot.Slot, 1)}},
	}
	p := &Player{c: c}

	container, err := p.OpenContainer(protocol.Position{1, 2, 3}, 0)
	if err == nil {
		t.Fatal("OpenContainer returned nil error for air block")
	}
	if container != nil {
		t.Fatalf("container = %#v, want nil", container)
	}
	if c.writes != 0 {
		t.Fatalf("WritePacket calls = %d, want 0", c.writes)
	}
}

type openContainerTestClient struct {
	world     bot.World
	inventory bot.InventoryHandler
	writes    int
}

func (c *openContainerTestClient) Connect(context.Context, string, *bot.ConnectOptions) error {
	return nil
}
func (c *openContainerTestClient) HandleGame(context.Context) error { return nil }
func (c *openContainerTestClient) Close(context.Context) error      { return nil }
func (c *openContainerTestClient) IsConnected() bool                { return true }
func (c *openContainerTestClient) WritePacket(context.Context, server.ServerboundPacket) error {
	c.writes++
	return nil
}
func (c *openContainerTestClient) PacketHandler() bot.PacketHandler {
	return openContainerTestPacketHandler{}
}
func (c *openContainerTestClient) EventHandler() bot.EventHandler {
	return openContainerTestEventHandler{}
}
func (c *openContainerTestClient) World() bot.World                { return c.world }
func (c *openContainerTestClient) Inventory() bot.InventoryHandler { return c.inventory }
func (c *openContainerTestClient) Player() bot.Player              { return nil }

type openContainerTestWorld struct {
	blk block.Block
	err error
}

func (w openContainerTestWorld) GetBlock(protocol.Position) (block.Block, error) {
	if w.err != nil {
		return nil, w.err
	}
	return w.blk, nil
}
func (w openContainerTestWorld) SetBlock(protocol.Position, block.Block) error { return nil }
func (w openContainerTestWorld) GetNearbyBlocks(protocol.Position, int32) ([]block.Block, error) {
	return nil, nil
}
func (w openContainerTestWorld) FindNearbyBlock(protocol.Position, int32, block.Block) (protocol.Position, error) {
	return protocol.Position{}, errors.New("not implemented")
}
func (w openContainerTestWorld) Entities() []bot.Entity               { return nil }
func (w openContainerTestWorld) GetEntity(int32) bot.Entity           { return nil }
func (w openContainerTestWorld) GetNearbyEntities(int32) []bot.Entity { return nil }
func (w openContainerTestWorld) GetEntitiesByType(entity.ID) []bot.Entity {
	return nil
}

type openContainerTestInventory struct {
	id        int32
	container bot.Container
}

func (i openContainerTestInventory) Inventory() bot.Container               { return i.container }
func (i openContainerTestInventory) Container() bot.Container               { return i.container }
func (i openContainerTestInventory) CurrentContainerID() int32              { return i.id }
func (i openContainerTestInventory) Click(int32, int16, int32, int32) error { return nil }
func (i openContainerTestInventory) Close()                                 {}

type openContainerTestContainer struct {
	slots []slot.Slot
}

func (c openContainerTestContainer) GetSlot(index int) slot.Slot { return c.slots[index] }
func (c openContainerTestContainer) Slots() []slot.Slot          { return c.slots }
func (c openContainerTestContainer) SlotCount() int              { return len(c.slots) }
func (c openContainerTestContainer) FindEmpty() int16            { return -1 }
func (c openContainerTestContainer) FindItem(item.ID) int16      { return -1 }
func (c openContainerTestContainer) Click(int16, int32, int32) error {
	return nil
}

type openContainerTestPacketHandler struct{}

func (openContainerTestPacketHandler) AddPacketHandler(packetid.ClientboundPacketID, func(context.Context, gameclient.ClientboundPacket)) {
}
func (openContainerTestPacketHandler) AddRawPacketHandler(packetid.ClientboundPacketID, func(context.Context, packet.Packet)) {
}
func (openContainerTestPacketHandler) AddGenericPacketHandler(func(context.Context, gameclient.ClientboundPacket)) {
}
func (openContainerTestPacketHandler) HandlePacket(context.Context, gameclient.ClientboundPacket) {}

type openContainerTestEventHandler struct{}

func (openContainerTestEventHandler) PublishEvent(string, any) error { return nil }
func (openContainerTestEventHandler) SubscribeEvent(string, func(any) error) {
}
