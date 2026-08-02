package main

import (
	"autocraft/config"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KonjacBot/go-mc/level/block"
	"github.com/KonjacBot/go-mc/level/item"
	"github.com/KonjacBot/minego/pkg/auth"
	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/client"
	"github.com/KonjacBot/minego/pkg/game/player"
	"github.com/KonjacBot/minego/pkg/protocol"
	cp "github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/pkg/protocol/slot/display/recipe"
	dislot "github.com/KonjacBot/minego/pkg/protocol/slot/display/slot"
	"github.com/go-gl/mathgl/mgl64"
)

var c bot.Client
var cfg config.Config
var glassRID int32
var inventoryInteraction sync.Mutex
var selectConfiguredServer sync.Once

const (
	normalCraftRetryDelay  = 5 * time.Millisecond
	blockedCraftRetryDelay = 10 * time.Millisecond
)

var startCraftLoop = sync.OnceFunc(func() {
	time.Sleep(500 * time.Millisecond)
	reporter := craftLoopErrorReporter{}
	for {
		err := craft()
		if reporter.shouldLog(err) {
			fmt.Println(err)
		}
		time.Sleep(craftRetryDelay(err))
	}
})

func main() {
	var err error
	cfg, err = config.ReadConfig()
	if err != nil {
		return
	}

	c = client.NewClient(&bot.ClientOptions{AuthProvider: &auth.KonjacAuth{
		UserCode: cfg.UserCode,
	}})

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	bot.SubscribeEvent(c, func(e player.MessageEvent) error {
		message := e.Message.ClearString()
		if message == loadedMessage {
			c.WritePacket(ctx, &server.ClientCommand{
				Action: 0,
			})
			go func() {
				startCraftLoop()
			}()
		}
		if balance, ok := parseEmeraldBalance(message); ok {
			emeraldBalance.Store(balance)
		}
		handlePrivateMessage(message)
		handleTeleportRequest(message)

		fmt.Println(e.Message.String())
		return nil
	})

	bot.AddHandler(c, func(ctx context.Context, p *cp.RecipeBookAdd) {
		for _, r := range p.Recipes {
			rID := r.RecipeID
			if r.Display.Display.RecipeType() == recipe.DisplayCraftingShaped {
				shaped := r.Display.Display.(*recipe.Shaped)
				switch s := shaped.Result.SlotDisplay.(type) {
				case *dislot.Item:
					if s.ID == int32(item.GlassPane{}.ID()) {
						glassRID = rID
					}

				case *dislot.ItemStack:
					if int32(s.ItemStack.ItemID) == int32(item.GlassPane{}.ID()) {
						glassRID = rID
					}
				}

			}
		}
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.SetHealth) {
		c.WritePacket(ctx, &server.ClientCommand{
			Action: 0,
		})
	})
	bot.AddHandler(c, func(context.Context, *cp.Login) {
		selectConfiguredServer.Do(func() {
			command, err := sendConfiguredServer(c, cfg.Server)
			if command == "" {
				fmt.Println("select configured server: server is empty")
				return
			}
			fmt.Printf("select configured server: %s\n", command)
			if err != nil {
				fmt.Printf("select configured server: %v\n", err)
			}
		})
	})
	bot.AddHandler(c, func(ctx context.Context, p *cp.SetTabListHeaderAndFooter) {
		message := p.Header.ClearString() + "\n" + p.Footer.ClearString()
		if balance, ok := parseEmeraldBalance(message); ok {
			emeraldBalance.Store(balance)
		}
	})

	err = c.Connect(ctx, cfg.Address, &bot.ConnectOptions{
		FakeHost: "mcfallout.net",
		Proxy:    cfg.Proxy,
	})
	if err != nil {
		panic(err)
	}

	err = c.HandleGame(ctx)
	if err != nil {
		panic(err)
	}
}

func craft() error {
	inventoryInteraction.Lock()
	defer inventoryInteraction.Unlock()

	runCraftWorkflow(craftWorkflow{
		craftGlass:   craftGlass,
		putGlassPane: putGlassPane,
		takeGlass:    takeGlass,
		discardJunk:  discardJunk,
	})

	return nil
}

type craftWorkflow struct {
	craftGlass   func() (int32, int32)
	putGlassPane func()
	takeGlass    func()
	discardJunk  func()
}

func runCraftWorkflow(workflow craftWorkflow) {
	glassCount, glassPaneCount := workflow.craftGlass()

	if glassPaneCount > 0 {
		workflow.putGlassPane()
	}

	if glassCount < 64 {
		workflow.takeGlass()
	}

	workflow.discardJunk()
}

func putGlassPane() {
	for _, pos := range cfg.PlacePos {
		container, err := openContainerWithoutPlacing(pos)
		if err != nil || container == nil {
			fmt.Println(err)
			continue
		}
		c.Player().CheckServer()

		for i, s := range container.Slots() {
			if i >= 27 && s.ItemID == (item.GlassPane{}.ID()) {
				_ = container.Click(int16(i), 1, 0)
				time.Sleep(5 * time.Millisecond)
			}
			if i < 27 && (s.ItemID != 0 && s.ItemID != (item.GlassPane{}.ID())) {
				_ = container.Click(int16(i), 1, 0)
				time.Sleep(5 * time.Millisecond)
			}
		}
		time.Sleep(15 * time.Millisecond)
	}
}

func takeGlass() {
	container, err := openContainerWithoutPlacing(cfg.TakePos)
	if err != nil || container == nil {
		return
	}
	c.Player().CheckServer()
	count := 0

	for i, s := range container.Slots() {
		if i < 27 && s.ItemID == (item.Glass{}.ID()) {
			_ = container.Click(int16(i), 1, 0)
			time.Sleep(15 * time.Millisecond)
			count++
		}
		if count > 14 {
			break
		}
	}

	f := sync.OnceFunc(func() {
		c.Player().Entity().SetRotation(mgl64.Vec2{c.Player().Entity().Rotation()[0], 0})
		c.Player().UpdateLocation()
		time.Sleep(500 * time.Millisecond)
		c.Player().CheckServer()
	})

	for i, s := range container.Slots() {
		if i < 27 && (s.ItemID != (item.Glass{}.ID()) && s.ItemID != 0) {
			f()
			_ = container.Click(int16(i), 4, 1)
			time.Sleep(5 * time.Millisecond)
		}
	}

	c.Player().Entity().SetRotation(mgl64.Vec2{c.Player().Entity().Rotation()[0], 90})
	c.Player().UpdateLocation()
	if count >= 0 {
		time.Sleep(500 * time.Millisecond)
	}
}

func craftGlass() (int32, int32) {
	playerPos := c.Player().Entity().Position()
	pos := protocol.Position{int32(playerPos[0]), int32(playerPos[1]), int32(playerPos[2])}

	craftingTablePos, err := c.World().FindNearbyBlock(pos, 6, block.CraftingTable{})
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	con, err := openContainerWithoutPlacing(craftingTablePos)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}

	c.Player().CheckServer()

	craftAllGlassPanes(con)
	glassCount := 0
	glassPaneCount := 0
	ff := false
	f := sync.OnceFunc(func() {
		c.Player().Entity().SetRotation(mgl64.Vec2{c.Player().Entity().Rotation()[0], 0})
		c.Player().UpdateLocation()
		time.Sleep(500 * time.Millisecond)
		c.Player().CheckServer()
	})

	for i, s := range con.Slots() {
		if s.ItemID == (item.Glass{}.ID()) {
			glassCount += int(s.Count)
			continue
		}
		if s.ItemID == (item.GlassPane{}.ID()) {
			glassPaneCount += int(s.Count)
			continue
		}
		if s.ItemID != 0 {
			f()
			ff = true

			_ = con.Click(int16(i), 4, 1)
			time.Sleep(5 * time.Millisecond)
		}
	}
	if ff {
		c.Player().Entity().SetRotation(mgl64.Vec2{c.Player().Entity().Rotation()[0], 90})
		c.Player().UpdateLocation()
	}
	return int32(glassCount), int32(glassPaneCount)
}

func craftAllGlassPanes(container bot.Container) {
	_ = c.WritePacket(context.Background(), &server.PlaceRecipe{WindowID: c.Inventory().CurrentContainerID(), RecipeID: glassRID, MakeAll: true})

	_ = container.Click(0, 1, 0)
	time.Sleep(5 * time.Millisecond)
}

func openContainerWithoutPlacing(pos protocol.Position) (bot.Container, error) {
	c.Inventory().Close()
	return c.Player().OpenContainer(pos, 0)
}

type craftLoopErrorReporter struct {
	lastError string
}

func (r *craftLoopErrorReporter) shouldLog(err error) bool {
	if err == nil {
		r.lastError = ""
		return false
	}

	message := err.Error()
	if message == r.lastError {
		return false
	}
	r.lastError = message
	return true
}

func craftRetryDelay(err error) time.Duration {
	if err != nil {
		return blockedCraftRetryDelay
	}
	return normalCraftRetryDelay
}
