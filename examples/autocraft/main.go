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

	f := sync.OnceFunc(func() {
		time.Sleep(500 * time.Millisecond)
		for {
			craft()
		}
	})

	bot.SubscribeEvent(c, func(e player.MessageEvent) error {
		message := e.Message.ClearString()
		if message == "[系統] 讀取人物成功。" {
			c.WritePacket(ctx, &server.ClientCommand{
				Action: 0,
			})
			go func() {
				f()
			}()
		}

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

func craft() {

	glassCount, glassPaneCount := craftGlass()

	if glassPaneCount > 0 {
		putGlassPane()
	}

	if glassCount < 64 {
		takeGlass()
	}

	fmt.Println(glassCount, glassPaneCount)

}

func putGlassPane() {
	for _, pos := range cfg.PlacePos {
		c.Inventory().Close()
		container, err := c.Player().OpenContainer(pos, 1)
		if err != nil || container == nil {
			fmt.Println(err)
			continue
		}
		c.Player().CheckServer()
		time.Sleep(500 * time.Millisecond)

		for i, s := range container.Slots() {
			if i >= 27 && s.ItemID == (item.GlassPane{}.ID()) {
				_ = container.Click(int16(i), 1, 0)
				time.Sleep(50 * time.Millisecond)
			}
			if i < 27 && (s.ItemID != 0 && s.ItemID != (item.GlassPane{}.ID())) {
				_ = container.Click(int16(i), 1, 0)
				time.Sleep(50 * time.Millisecond)
			}
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func takeGlass() {
	c.Inventory().Close()
	container, err := c.Player().OpenContainer(cfg.TakePos, 1)
	if err != nil || container == nil {
		return
	}
	c.Player().CheckServer()
	count := 0

	for i, s := range container.Slots() {
		if i < 27 && s.ItemID == (item.Glass{}.ID()) {
			_ = container.Click(int16(i), 1, 0)
			time.Sleep(50 * time.Millisecond)
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
			time.Sleep(50 * time.Millisecond)
		}
	}

	c.Player().Entity().SetRotation(mgl64.Vec2{c.Player().Entity().Rotation()[0], 90})
	c.Player().UpdateLocation()
	if count >= 0 {
		time.Sleep(500 * time.Millisecond)
	}
}

func craftGlass() (int32, int32) {
	c.Inventory().Close()

	playerPos := c.Player().Entity().Position()
	pos := protocol.Position{int32(playerPos[0]), int32(playerPos[1]), int32(playerPos[2])}

	craftingTablePos, err := c.World().FindNearbyBlock(pos, 6, block.CraftingTable{})
	if err != nil {
		return 0, 0
	}
	con, err := c.Player().OpenContainer(craftingTablePos, 1)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}

	c.Player().CheckServer()

	for i := 0; i < 6; i++ {
		_ = c.WritePacket(context.Background(), &server.PlaceRecipe{WindowID: c.Inventory().CurrentContainerID(), RecipeID: glassRID, MakeAll: true})

		_ = con.Click(0, 1, 0)
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(150 * time.Millisecond)
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
			time.Sleep(50 * time.Millisecond)
		}
	}
	if ff {
		c.Player().Entity().SetRotation(mgl64.Vec2{c.Player().Entity().Rotation()[0], 90})
		c.Player().UpdateLocation()
	}
	return int32(glassCount), int32(glassPaneCount)
}
