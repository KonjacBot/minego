package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KonjacBot/go-mc/level/item"
	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot"
)

func Test_parseEmeraldBalance_whenTabListContainsBalance(t *testing.T) {
	balance, ok := parseEmeraldBalance("綠寶石餘額 : 12,345 / 村民錠餘額 : 0")

	if !ok {
		t.Fatal("expected balance match")
	}
	if balance != 12345 {
		t.Fatalf("balance = %d, want 12345", balance)
	}
}

func Test_teleportPlayer_whenTeleportToOwnerMessage(t *testing.T) {
	playerName, ok := teleportPlayer("[系統] PatyHank 想要你傳送到 該玩家 的位置")

	if !ok {
		t.Fatal("expected teleport match")
	}
	if playerName != "PatyHank" {
		t.Fatalf("playerName = %q, want PatyHank", playerName)
	}
}

func TestOpenContainerWithoutPlacingWhenHotbarIsFull(t *testing.T) {
	originalClient := c
	t.Cleanup(func() { c = originalClient })

	inventory := &openingTestInventory{}
	player := &openingTestPlayer{container: &openingTestContainer{}}
	c = &openingTestClient{inventory: inventory, player: player}
	position := protocol.Position{12, 64, -8}

	container, err := openContainerWithoutPlacing(position)
	if err != nil {
		t.Fatalf("openContainerWithoutPlacing() error = %v, want nil", err)
	}
	if !inventory.closed {
		t.Fatal("expected the previous container to be closed")
	}
	if !player.opened {
		t.Fatal("expected a container open attempt despite the full hotbar")
	}
	if player.position != position || player.hand != 0 {
		t.Fatalf("OpenContainer() = (%v, %d), want (%v, 0)", player.position, player.hand, position)
	}
	if container != player.container {
		t.Fatal("expected the container returned by Player.OpenContainer")
	}
}

func TestRunCraftWorkflowDiscardsJunk(t *testing.T) {
	garbageDiscarded := false

	runCraftWorkflow(craftWorkflow{
		craftGlass: func() (int32, int32) { return 64, 0 },
		putGlassPane: func() {
			t.Fatal("did not expect glass panes to be stored")
		},
		takeGlass: func() {
			t.Fatal("did not expect glass to be taken")
		},
		discardJunk: func() { garbageDiscarded = true },
	})

	if !garbageDiscarded {
		t.Fatal("expected every craft workflow to discard backpack junk")
	}
}

func TestCraftAllGlassPanesUsesOneMakeAllRequest(t *testing.T) {
	originalClient, originalRecipeID := c, glassRID
	t.Cleanup(func() {
		c = originalClient
		glassRID = originalRecipeID
	})

	client := &openingTestClient{inventory: recipeTestInventory{}}
	c = client
	glassRID = 42
	container := &recipeTestContainer{}

	craftAllGlassPanes(container)

	if len(client.writes) != 1 {
		t.Fatalf("MakeAll requests = %d, want 1", len(client.writes))
	}
	if len(container.clicks) != 1 {
		t.Fatalf("result clicks = %d, want 1", len(container.clicks))
	}
	request, ok := client.writes[0].(*server.PlaceRecipe)
	if !ok || request.RecipeID != glassRID || !request.MakeAll {
		t.Fatalf("recipe request = %#v, want MakeAll for recipe %d", client.writes[0], glassRID)
	}
}

func TestJunkInventorySlotsKeepsCraftMaterialsAndHotbar(t *testing.T) {
	inventory := &junkTestContainer{slots: make([]slot.Slot, 46)}
	inventory.slots[9] = slot.Slot{ItemID: item.Glass{}.ID(), Count: 64}
	inventory.slots[10] = slot.Slot{ItemID: item.GlassPane{}.ID(), Count: 64}
	inventory.slots[11] = slot.Slot{ItemID: item.Glass{}.ID() + 1_000, Count: 1}
	inventory.slots[36] = slot.Slot{ItemID: item.Glass{}.ID() + 1_000, Count: 1}

	got := junkInventorySlots(inventory)
	if len(got) != 1 || got[0] != 11 {
		t.Fatalf("junkInventorySlots() = %v, want [11]", got)
	}
}

func TestServerCommand(t *testing.T) {
	command, ok := serverCommand(" server71 ")
	if !ok || command != "server server71" {
		t.Fatalf("serverCommand() = (%q, %t), want (\"server server71\", true)", command, ok)
	}

	if command, ok := serverCommand(" "); ok || command != "" {
		t.Fatalf("serverCommand(empty) = (%q, %t), want (\"\", false)", command, ok)
	}
}

func TestSendConfiguredServerCommandsPlayerDirectly(t *testing.T) {
	player := &openingTestPlayer{}
	client := &openingTestClient{player: player}

	command, err := sendConfiguredServer(client, "server71")
	if err != nil {
		t.Fatal(err)
	}
	if command != "server server71" || player.command != command {
		t.Fatalf("command = %q, player command = %q, want server server71", command, player.command)
	}
}

func TestCraftLoopErrorReporterSuppressesRepeatedErrors(t *testing.T) {
	reporter := craftLoopErrorReporter{}
	fullHotbar := errors.New("no empty hotbar slot")

	if !reporter.shouldLog(fullHotbar) {
		t.Fatal("expected the first error to be logged")
	}
	if reporter.shouldLog(fullHotbar) {
		t.Fatal("expected a repeated error to be suppressed")
	}
	if reporter.shouldLog(nil) {
		t.Fatal("expected a successful iteration not to be logged")
	}
	if !reporter.shouldLog(fullHotbar) {
		t.Fatal("expected an error after a successful iteration to be logged")
	}
}

func TestCraftRetryDelay(t *testing.T) {
	if got := craftRetryDelay(nil); got != normalCraftRetryDelay {
		t.Fatalf("craftRetryDelay(nil) = %s, want %s", got, normalCraftRetryDelay)
	}
	if got := craftRetryDelay(errors.New("open container failed")); got != blockedCraftRetryDelay {
		t.Fatalf("craftRetryDelay(err) = %s, want %s", got, blockedCraftRetryDelay)
	}
	if blockedCraftRetryDelay <= time.Second {
		t.Fatalf("blocked retry delay = %s, want more than one second", blockedCraftRetryDelay)
	}
}

type openingTestClient struct {
	bot.Client
	inventory bot.InventoryHandler
	player    bot.Player
	writes    []server.ServerboundPacket
}

func (c *openingTestClient) Inventory() bot.InventoryHandler { return c.inventory }
func (c *openingTestClient) Player() bot.Player              { return c.player }
func (c *openingTestClient) WritePacket(_ context.Context, packet server.ServerboundPacket) error {
	c.writes = append(c.writes, packet)
	return nil
}

type openingTestInventory struct {
	bot.InventoryHandler
	closed bool
}

func (i *openingTestInventory) Close() { i.closed = true }
func (*openingTestInventory) Inventory() bot.Container {
	return fullHotbarTestContainer{}
}

type fullHotbarTestContainer struct{ bot.Container }

func (fullHotbarTestContainer) GetSlot(int) slot.Slot { return slot.Slot{Count: 64} }
func (fullHotbarTestContainer) SlotCount() int        { return 46 }

type openingTestPlayer struct {
	bot.Player
	container bot.Container
	opened    bool
	position  protocol.Position
	hand      int32
	command   string
}

func (p *openingTestPlayer) Command(command string) error {
	p.command = command
	return nil
}

func (p *openingTestPlayer) OpenContainer(pos protocol.Position, hand int32) (bot.Container, error) {
	p.opened = true
	p.position = pos
	p.hand = hand
	return p.container, nil
}

type openingTestContainer struct{ bot.Container }

type recipeTestInventory struct{ bot.InventoryHandler }

func (recipeTestInventory) CurrentContainerID() int32 { return 7 }

type recipeTestContainer struct {
	bot.Container
	clicks []int16
}

func (c *recipeTestContainer) Click(index int16, _ int32, _ int32) error {
	c.clicks = append(c.clicks, index)
	return nil
}

type junkTestContainer struct {
	bot.Container
	slots []slot.Slot
}

func (c *junkTestContainer) GetSlot(index int) slot.Slot { return c.slots[index] }
func (c *junkTestContainer) SlotCount() int              { return len(c.slots) }
