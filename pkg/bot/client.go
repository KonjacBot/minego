package bot

import (
	"context"

	"git.konjactw.dev/patyhank/minego/pkg/auth"
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/server"
)

type Client interface {
	Connect(ctx context.Context, addr string, options *ConnectOptions) error
	HandleGame(ctx context.Context) error
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

type ProxyConfig struct {
	Type     string `json:"type" toml:"type"`
	Host     string `json:"host" toml:"host"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
}

type ConnectOptions struct {
	FakeHost string       `json:"fake_host,omitempty" toml:"fake_host,omitempty"`
	Proxy    *ProxyConfig `json:"proxy,omitempty" toml:"proxy,omitempty"`
}
