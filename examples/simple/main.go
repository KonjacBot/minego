package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/KonjacBot/minego/v2/pkg/auth"
	"github.com/KonjacBot/minego/v2/pkg/bot"
	"github.com/KonjacBot/minego/v2/pkg/client"
	"github.com/KonjacBot/minego/v2/pkg/game/player"
)

func main() {
	userCode := os.Getenv("MINEGO_USER_CODE")
	if userCode == "" {
		log.Fatal("MINEGO_USER_CODE is required")
	}
	address := os.Getenv("MINEGO_SERVER")
	if address == "" {
		address = "localhost:25565"
	}

	c := client.NewClient(&bot.ClientOptions{AuthProvider: &auth.KonjacAuth{
		UserCode: userCode,
	}})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := c.Connect(ctx, address, nil); err != nil {
		log.Fatal(err)
	}
	defer c.Close(context.Background())

	c.On(func(_ context.Context, e player.MessageEvent) {
		fmt.Println(e.Message.String())
	})

	if err := c.HandleGame(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
