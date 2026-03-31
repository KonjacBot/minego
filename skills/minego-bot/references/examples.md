# minego Example Bots

## Simple Bot

A minimal bot that connects to a server and prints chat messages.

```go
package main

import (
	"context"
	"fmt"

	"github.com/KonjacBot/minego/pkg/auth"
	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/client"
	"github.com/KonjacBot/minego/pkg/game/player"
)

func main() {
	userCode := "your-code"
	c := client.NewClient(&bot.ClientOptions{AuthProvider: &auth.KonjacAuth{
		UserCode: userCode,
	}})

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err := c.Connect(ctx, "mcfallout.net", nil)
	if err != nil {
		panic(err)
	}

	bot.SubscribeEvent(c, func(e player.MessageEvent) error {
		fmt.Println(e.Message.String())
		return nil
	})

	err = c.HandleGame(ctx)
	if err != nil {
		panic(err)
	}
}
```

**Key points:**
- `client.NewClient()` creates the client with auth options
- `Connect()` establishes the connection (pass `nil` for default options)
- `bot.SubscribeEvent()` subscribes to typed events before entering the game loop
- `HandleGame()` is the blocking main loop — call it last

---

## Autocraft Bot

An advanced bot that automates crafting glass panes from glass blocks using a crafting table and chest containers.

### Configuration (TOML)

```go
type Config struct {
	Address     string              `toml:"address"`
	Proxy       *bot.ProxyConfig    `toml:"proxy,omitempty"`
	UserCode    string              `toml:"user_code"`
	TakePos     protocol.Position   `toml:"take_pos"`
	TakeButton  protocol.Position   `toml:"take_button"`
	PlacePos    []protocol.Position `toml:"place_pos"`
	PlaceButton []protocol.Position `toml:"place_button"`
}
```

### Bot Setup

```go
c = client.NewClient(&bot.ClientOptions{AuthProvider: &auth.KonjacAuth{
	UserCode: cfg.UserCode,
}})

// Subscribe to chat messages to detect login
bot.SubscribeEvent(c, func(e player.MessageEvent) error {
	message := e.Message.ClearString()
	if message == "[系統] 讀取人物成功。" {
		// Trigger actions after login
		c.WritePacket(ctx, &server.ClientCommand{Action: 0})
		go startCrafting()
	}
	fmt.Println(e.Message.String())
	return nil
})

// Handle recipe discovery
bot.AddHandler(c, func(ctx context.Context, p *cp.RecipeBookAdd) {
	for _, r := range p.Recipes {
		// Find specific recipe IDs at runtime
	}
})

// Auto-respawn on death
bot.AddHandler(c, func(ctx context.Context, p *cp.SetHealth) {
	c.WritePacket(ctx, &server.ClientCommand{Action: 0})
})

err = c.Connect(ctx, cfg.Address, &bot.ConnectOptions{
	FakeHost: "mcfallout.net",
	Proxy:    cfg.Proxy,
})
c.HandleGame(ctx)
```

### Container Interaction Pattern

```go
// Open a container at a block position
container, err := c.Player().OpenContainer(pos, 1)
if err != nil || container == nil {
	return
}
c.Player().CheckServer()
time.Sleep(500 * time.Millisecond)

// Iterate slots and shift-click items
for i, s := range container.Slots() {
	if i >= 27 && s.ItemID == item.GlassPane{}.ID() {
		_ = container.Click(int16(i), 1, 0)  // shift-click
		time.Sleep(50 * time.Millisecond)
	}
}
```

### Finding and Using a Crafting Table

```go
playerPos := c.Player().Entity().Position()
pos := protocol.Position{int32(playerPos[0]), int32(playerPos[1]), int32(playerPos[2])}

craftingTablePos, err := c.World().FindNearbyBlock(pos, 6, block.CraftingTable{})
if err != nil {
	return
}

con, err := c.Player().OpenContainer(craftingTablePos, 1)
if err != nil {
	return
}
c.Player().CheckServer()

// Place recipe in crafting table
c.WritePacket(ctx, &server.PlaceRecipe{
	WindowID: c.Inventory().CurrentContainerID(),
	RecipeID: glassRID,
	MakeAll:  true,
})

// Take result from slot 0
con.Click(0, 1, 0)
```

### Dropping Unwanted Items

```go
// Rotate player to face drop direction
c.Player().Entity().SetRotation(mgl64.Vec2{yaw, 0})
c.Player().UpdateLocation()
time.Sleep(500 * time.Millisecond)
c.Player().CheckServer()

// Drop full stack from a container slot
con.Click(int16(slotIndex), 4, 1)
```

**Key points:**
- Always call `CheckServer()` after opening containers to sync state
- Add `time.Sleep()` between click operations for reliability
- Use `sync.OnceFunc` to ensure one-time setup actions
- Close containers with `c.Inventory().Close()` before opening new ones
- Item IDs use zero-value struct pattern: `item.Glass{}.ID()`
- Block types use zero-value struct pattern: `block.CraftingTable{}`
