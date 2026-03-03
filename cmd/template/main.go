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
	userCode := "powru"
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
