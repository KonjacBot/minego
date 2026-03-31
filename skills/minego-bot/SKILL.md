---
name: minego-bot
description: "Build Minecraft bots using the minego Go library. Use when: writing Go code with minego, creating Minecraft bots, connecting to Minecraft servers, handling packets, subscribing to game events, managing inventory, player movement, world interaction, pathfinding, crafting automation."
argument-hint: "Describe the Minecraft bot behavior you want to build"
---

# minego Bot Development Guide

## When to Use

- Creating a new Minecraft bot in Go using `github.com/KonjacBot/minego`
- Connecting to a Minecraft Java server, handling game events, or automating tasks
- Working with inventory, player movement, world blocks, entities, or crafting
- Extending an existing minego bot with new features

## Module

```
github.com/KonjacBot/minego
```

## Core Architecture

The library is split into **interface** (`pkg/bot/`) and **implementation** (`pkg/client/`, `pkg/game/`) layers:

| Package | Purpose |
|---------|---------|
| `pkg/bot` | Interfaces: `Client`, `Player`, `World`, `Entity`, `InventoryHandler`, `Container` |
| `pkg/client` | `NewClient()` constructor and connection logic |
| `pkg/auth` | Authentication providers (`auth.Provider` interface) |
| `pkg/game/player` | Player movement, chat, pathfinding, events |
| `pkg/game/world` | Block/entity access, chunk management |
| `pkg/game/inventory` | Container/slot management, click operations |
| `pkg/protocol` | `Position`, codecs, packet definitions |
| `pkg/protocol/packet/game/server` | Serverbound packets (client → server) |
| `pkg/protocol/packet/game/client` | Clientbound packets (server → client) |

## Step-by-Step: Creating a Bot

### 1. Initialize module and import

```go
import (
    "context"
    "github.com/KonjacBot/minego/pkg/auth"
    "github.com/KonjacBot/minego/pkg/bot"
    "github.com/KonjacBot/minego/pkg/client"
)
```

### 2. Create and connect the client

```go
c := client.NewClient(&bot.ClientOptions{
    AuthProvider: &auth.KonjacAuth{UserCode: "your-code"},
})

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := c.Connect(ctx, "server-address:25565", &bot.ConnectOptions{
    FakeHost: "optional.hostname",  // optional
    Proxy:    nil,                  // or &bot.ProxyConfig{Type: "socks5", Host: "..."}
})
if err != nil {
    panic(err)
}
```

### 3. Register event/packet handlers BEFORE or AFTER Connect, but BEFORE HandleGame

Register events and packet handlers. Then start the blocking game loop:

```go
// Subscribe to typed events
bot.SubscribeEvent(c, func(e player.MessageEvent) error {
    fmt.Println(e.Message.String())
    return nil
})

// Add typed packet handler (for clientbound packets)
bot.AddHandler(c, func(ctx context.Context, p *cp.SetHealth) {
    // respond to health changes
})

// Start the blocking game loop (must be last)
err = c.HandleGame(ctx)
```

### 4. Send packets to the server

```go
c.WritePacket(ctx, &server.ClientCommand{Action: 0})
```

## Key Interfaces

### `bot.Client`

```go
Connect(ctx, addr, *ConnectOptions) error
HandleGame(ctx) error
Close(ctx) error
IsConnected() bool
WritePacket(ctx, ServerboundPacket) error
Player() Player
World() World
Inventory() InventoryHandler
PacketHandler() PacketHandler
EventHandler() EventHandler
```

### `bot.Player`

```go
Entity() Entity                    // get player's entity (position, rotation)
FlyTo(pos mgl64.Vec3) error
WalkTo(pos mgl64.Vec3) error
LookAt(pos mgl64.Vec3) error
UpdateLocation()                   // send current position/rotation to server
BreakBlock(pos protocol.Position) error
PlaceBlock(pos protocol.Position) error
OpenContainer(pos protocol.Position, hand int32) (Container, error)
OpenMenu(command string) (Container, error)
Command(command string) error
Chat(message string) error
CheckServer()                      // sync state with server
UseItem(hand int8) error
```

### `bot.World`

```go
GetBlock(pos protocol.Position) (block.Block, error)
SetBlock(pos protocol.Position, b block.Block) error
FindNearbyBlock(pos protocol.Position, radius int32, blk block.Block) (Position, error)
GetNearbyBlocks(pos protocol.Position, radius int32) ([]block.Block, error)
Entities() []Entity
GetEntity(id int32) Entity
GetNearbyEntities(radius int32) []Entity
GetEntitiesByType(entityType entity.ID) []Entity
```

### `bot.InventoryHandler`

```go
Inventory() Container       // player inventory
Container() Container       // currently open container
CurrentContainerID() int32
Click(container, slot int16, mode, button int32) error
Close()
```

### `bot.Container`

```go
Slots() []slot.Slot
GetSlot(index int) slot.Slot
SlotCount() int
FindEmpty() int16
FindItem(itemID item.ID) int16
Click(slot int16, mode int32, button int32) error
```

### `bot.Entity`

```go
ID() int32
UUID() uuid.UUID
Type() entity.ID
Position() mgl64.Vec3
Rotation() mgl64.Vec2
SetPosition(pos mgl64.Vec3)
SetRotation(rot mgl64.Vec2)
Metadata() map[uint8]metadata.Metadata
Equipment() map[int8]slot.Slot
```

### `protocol.Position`

```go
type Position [3]int32  // [X, Y, Z]
// Methods: DistanceTo, Add, Sub, Mul, Div, IsZero, Clone, Equals
```

## Event System

Subscribe to typed events with `bot.SubscribeEvent`:

| Event | ID | Fields |
|-------|----|--------|
| `player.MessageEvent` | `player:message` | `Message chat.Message` |
| `inventory.ContainerOpenEvent` | `inventory:container_open` | `WindowID`, `Type`, `Title` |
| `world.EntityAddEvent` | `world:entity_add` | entity data |
| `world.EntityRemoveEvent` | `world:entity_remove` | entity data |

```go
bot.SubscribeEvent(c, func(e player.MessageEvent) error {
    fmt.Println(e.Message.String())
    return nil
})
```

## Packet Handler System

Handle specific clientbound packets with `bot.AddHandler`:

```go
import cp "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"

bot.AddHandler(c, func(ctx context.Context, p *cp.RecipeBookAdd) {
    // process recipes
})
```

Send serverbound packets with `WritePacket`:

```go
import "github.com/KonjacBot/minego/pkg/protocol/packet/game/server"

c.WritePacket(ctx, &server.ClientCommand{Action: 0})
c.WritePacket(ctx, &server.PlaceRecipe{
    WindowID: c.Inventory().CurrentContainerID(),
    RecipeID: recipeID,
    MakeAll:  true,
})
```

## Common Patterns

### Opening and interacting with containers

```go
container, err := c.Player().OpenContainer(pos, 1)  // hand=1
if err != nil || container == nil {
    return
}
c.Player().CheckServer()
time.Sleep(500 * time.Millisecond)  // wait for server sync

for i, s := range container.Slots() {
    if s.ItemID == item.Glass{}.ID() {
        _ = container.Click(int16(i), 1, 0)  // shift-click
        time.Sleep(50 * time.Millisecond)
    }
}
```

### Finding nearby blocks

```go
playerPos := c.Player().Entity().Position()
pos := protocol.Position{int32(playerPos[0]), int32(playerPos[1]), int32(playerPos[2])}
tablePos, err := c.World().FindNearbyBlock(pos, 6, block.CraftingTable{})
```

### Player rotation and looking

```go
c.Player().Entity().SetRotation(mgl64.Vec2{yaw, pitch})
c.Player().UpdateLocation()
```

### Pathfinding

```go
import "github.com/KonjacBot/minego/pkg/game/player"

path, err := player.AStar(c.World(), startPos, goalPos)
```

### Proxy configuration

```go
&bot.ConnectOptions{
    Proxy: &bot.ProxyConfig{
        Type:     "socks5",
        Host:     "proxy-host:1080",
        Username: "user",
        Password: "pass",
    },
}
```

## Click Modes (Inventory)

| Mode | Button | Action |
|------|--------|--------|
| 0 | 0 | Left click (pick up / place full stack) |
| 0 | 1 | Right click (pick up / place single) |
| 1 | 0 | Shift + left click (quick move) |
| 1 | 1 | Shift + right click |
| 4 | 0 | Drop one item |
| 4 | 1 | Drop full stack |

## Reference Examples

See [simple example](./references/examples.md#simple-bot) and [autocraft example](./references/examples.md#autocraft-bot) for complete working bots.
