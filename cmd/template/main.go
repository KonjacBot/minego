package main

import (
	"context"
	"fmt"

	"git.konjactw.dev/patyhank/minego/pkg/auth"
	"git.konjactw.dev/patyhank/minego/pkg/bot"
	"git.konjactw.dev/patyhank/minego/pkg/client"
	"git.konjactw.dev/patyhank/minego/pkg/game/player"
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
