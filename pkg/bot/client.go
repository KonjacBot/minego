package bot

import (
	"context"

	"git.konjactw.dev/patyhank/minego/pkg/auth"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
)

type Client interface {
	Connect(ctx context.Context, addr string, options *ConnectOptions) error
	Close(ctx context.Context) error
	IsConnected() bool
	WritePacket(ctx context.Context, packet server.ServerboundPacket) error

	PacketHandler() PacketHandler
	EventHandler() EventHandler
	World() World
	Inventory() InventoryHandler
	Player() Player
}

type ClientOptions struct {
	AuthProvider auth.Provider
}

type ConnectOptions struct {
	FakeHost string
}
