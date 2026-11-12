---
name: minego-bot
description: "Build Minecraft bots using the minego Go library. Use when: writing Go code with minego, creating Minecraft bots, connecting to Minecraft servers, handling packets, subscribing to game events, managing inventory, player movement, world interaction, pathfinding, crafting automation."
argument-hint: "Describe the Minecraft bot behavior you want to build"
---

# minego Bot Development Guide

## When to Use

- Creating a new Minecraft bot in Go using `github.com/KonjacBot/minego/v2`
- Connecting to a Minecraft Java server, handling game events, or automating tasks
- Working with inventory, player movement, world blocks, entities, or crafting
- Extending an existing minego bot with new features

## Module

```
github.com/KonjacBot/minego/v2
```

## Core Architecture

The library is split into **interface** (`pkg/bot/`) and **implementation** (`pkg/client/`, `pkg/game/`) layers:

| Package | Purpose |
|---------|---------|
| `pkg/bot` | Interfaces: `Client`, `Player`, `World`, `Entity`, `InventoryHandler`, `Container` |
| `pkg/client` | `NewClient()` returns `*client.Client`: connection logic plus `On`/`WaitFor`/`Stream`/`Emit` |
| `pkg/bus` | The typed packet and event bus behind `On`, `OnCtx`, `Once`, `WaitFor`, `Stream`, `Emit`, `OnFrame`, `OnRaw` |
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
    "os"

    "github.com/KonjacBot/minego/v2/pkg/auth"
    "github.com/KonjacBot/minego/v2/pkg/bot"
    "github.com/KonjacBot/minego/v2/pkg/client"
)
```

### 2. Create and connect the client

```go
c := client.NewClient(&bot.ClientOptions{
    AuthProvider: &auth.KonjacAuth{UserCode: os.Getenv("MINEGO_USER_CODE")},
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
defer c.Close(context.Background())
```

### 3. Register event/packet handlers BEFORE or AFTER Connect, but BEFORE HandleGame

Register events and packet handlers. Then start the blocking game loop:

```go
// Subscribe to an event: the type comes from the handler's parameter
c.On(func(_ context.Context, e player.MessageEvent) {
    fmt.Println(e.Message.String())
})

// Subscribe to a clientbound packet the same way
c.On(func(ctx context.Context, p *cp.SetHealth) {
    // respond to health changes
})

// Start the blocking game loop (must be last)
err = c.HandleGame(ctx)
```

## Runtime Safety Rules

- Reuse one `Client`; do not call `Connect`, `HandleGame`, or `Close` concurrently.
- Register packet handlers before `HandleGame`. Registration is concurrency-safe, but deterministic startup avoids missing early packets.
- `WritePacket` serializes concurrent writes and observes an existing context deadline. Pass the handler context instead of `context.Background()` when available.
- All handlers (packets, events, `OnRaw`, `OnFrame`) run synchronously on the packet loop, in registration order. Keep them short; move long-running work to a goroutine that observes cancellation. Dispatches slower than 50 ms are logged.
- Every subscription returns `off func()`. Use `c.OnCtx(ctx, ...)` for handlers tied to a task so they are removed when the task ends; handlers that are never removed are a memory leak.
- A packet is decoded only if a typed handler subscribed to it. `OnRaw` payloads are only valid until the handler returns; copy what you keep.
- A panicking handler is logged and skipped; handlers do not return errors.
- `World`, `Entity`, `InventoryHandler`, and `Container` snapshots are safe to read while packets update state. `Container.Slots`, `Entity.Metadata`, and `Entity.Equipment` return copies.
- A missing chunk or block is unknown, not air. Pathfinding intentionally refuses routes through unloaded chunks.
- `OpenContainer` waits for the matching window and initial content. Still handle its error and a possible nil result.
- Interaction sequence IDs are managed by `Player`; do not construct interaction packets with inventory `StateID` as their sequence.

### 4. Send packets to the server

```go
if err := c.WritePacket(ctx, &server.ClientCommand{Action: 0}); err != nil {
    return err
}
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
Bus() *bus.Bus // interfaces cannot have generic methods: c.Bus().On(...)
```

`*client.Client` (what `client.NewClient` returns) implements `bot.Client` and
forwards the bus methods, so with the concrete type write `c.On(...)`:

```go
c.On(func(ctx context.Context, v T)) (off func())          // T = packet type or event type
c.OnCtx(ctx, func(ctx context.Context, v T)) (off func())  // removed when ctx ends
c.Once(func(ctx context.Context, v T)) (off func())
c.WaitFor(ctx, func(v T) bool) (T, error)                  // nil matcher = first one
c.Stream[T](ctx) iter.Seq[T]                                // for v := range c.Stream[T](ctx)
c.Emit(ctx, v)                                             // publish your own event types
c.OnFrame(func(id packetid.ClientboundPacketID, size int)) // every packet, no decoding
c.OnRaw(id, func(ctx context.Context, p pk.Packet))        // undecoded payload
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

Events are plain Go types; subscribe with `c.On` and the event type as the
handler's parameter:

| Event | Fields |
|-------|--------|
| `player.MessageEvent` | `Message chat.Message` |
| `inventory.ContainerOpenEvent` | `WindowID`, `Type`, `Title` |
| `world.EntityAddEvent` | `EntityID` |
| `world.EntityRemoveEvent` | `Entity` |
| `client.ConnectionStateChangeEvent` | `From`, `To` |

```go
c.On(func(_ context.Context, e player.MessageEvent) {
    fmt.Println(e.Message.String())
})

// Wait for something with a timeout instead of polling
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
e, err := c.WaitFor(ctx, func(e inventory.ContainerOpenEvent) bool { return e.WindowID > 0 })

// Your own event types work too
type Arrived struct{ Pos protocol.Position }
c.Emit(ctx, Arrived{Pos: pos})
```

## Packet Handler System

Handle specific clientbound packets with `c.On`; the packet type is the
handler's parameter:

```go
import cp "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"

c.On(func(ctx context.Context, p *cp.RecipeBookAdd) {
    // process recipes
})
```

Send serverbound packets with `WritePacket`:

```go
import "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"

c.WritePacket(ctx, &server.ClientCommand{Action: 0})
c.WritePacket(ctx, &server.PlaceRecipe{
    WindowID: c.Inventory().CurrentContainerID(),
    RecipeID: recipeID,
    MakeAll:  true,
})
```

Always check the returned errors in production code. They report cancellation, a closed connection, and packet write failures.

## Common Patterns

### Opening and interacting with containers

```go
container, err := c.Player().OpenContainer(pos, 1)  // hand=1
if err != nil || container == nil {
    return
}

for i, s := range container.Slots() {
    if s.ItemID == item.Glass{}.ID() {
		if err := container.Click(int16(i), 1, 0); err != nil {
			return
		}
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
import "github.com/KonjacBot/minego/v2/pkg/game/player"

path, err := player.AStar(c.World(), startPos, goalPos, 4096)
```

`AStar` returns `player.ErrMaxNodesExceeded` with the best partial path when it reaches the node limit. Decide explicitly whether the bot should follow that partial path.

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
