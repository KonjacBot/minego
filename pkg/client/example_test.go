package client_test

import (
	"context"
	"fmt"
	"time"

	"github.com/KonjacBot/minego/v2/pkg/bot"
	"github.com/KonjacBot/minego/v2/pkg/client"
	"github.com/KonjacBot/minego/v2/pkg/game/inventory"
	"github.com/KonjacBot/minego/v2/pkg/game/player"
	cp "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
)

// Subscribe to events and packets; the handler's parameter type selects what.
func ExampleClient_On() {
	c := client.NewClient(nil)

	c.On(func(_ context.Context, e player.MessageEvent) {
		fmt.Println(e.Message.String())
	})
	c.On(func(ctx context.Context, p *cp.SetHealth) {
		if p.Health <= 0 {
			_ = c.WritePacket(ctx, &server.ClientCommand{Action: 0})
		}
	})
}

// Wait for a condition instead of polling a flag set by a handler.
func ExampleClient_WaitFor() {
	c := client.NewClient(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opened, err := c.WaitFor(ctx, func(e inventory.ContainerOpenEvent) bool { return e.WindowID > 0 })
	_, _ = opened, err
}

// Tie handlers to a task so they are removed when it ends.
func ExampleClient_OnCtx() {
	c := client.NewClient(nil)

	task, done := context.WithCancel(context.Background())
	defer done()
	c.OnCtx(task, func(_ context.Context, p *cp.SystemChatMessage) {})
}

// Range over packets or events as they arrive.
func ExampleClient_Stream() {
	c := client.NewClient(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for e := range c.Stream[player.MessageEvent](ctx) {
			fmt.Println(e.Message.String())
		}
	}()
}

// Code that only holds the bot.Client interface subscribes through Bus.
func ExampleClient_Bus() {
	var c bot.Client = client.NewClient(nil)

	c.Bus().On(func(_ context.Context, p *cp.KeepAlive) {})
}
