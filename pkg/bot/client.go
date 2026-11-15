package bot

import (
	"context"
	"time"

	"github.com/KonjacBot/minego/v2/pkg/auth"
	"github.com/KonjacBot/minego/v2/pkg/bus"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
)

type Client interface {
	Connect(ctx context.Context, addr string, options *ConnectOptions) error
	HandleGame(ctx context.Context) error
	Close(ctx context.Context) error
	IsConnected() bool
	WritePacket(ctx context.Context, packet server.ServerboundPacket) error

	// Bus delivers packets and events. Interfaces cannot have generic
	// methods, so code holding a bot.Client subscribes through it:
	// c.Bus().On(func(ctx context.Context, p *client.KeepAlive) { ... }).
	Bus() *bus.Bus
	World() World
	Inventory() InventoryHandler
	Player() Player
}

type ClientOptions struct {
	AuthProvider auth.Provider
	// ReadIdleTimeout is renewed before each login, configuration, and play
	// packet read. Values at or below zero use the client default.
	ReadIdleTimeout time.Duration
	// ResourcePackPolicy controls the headless response to server resource
	// packs. The accept policy acknowledges successful loading without
	// downloading visual assets; decline reports that the pack was refused.
	ResourcePackPolicy ResourcePackPolicy
	// NoChunks stops the built-in World from storing chunks, for callers that
	// keep terrain in their own structure via packet handlers.
	NoChunks bool
	// ChunkLight keeps the light data of loaded chunks. Off by default: it is
	// about half of each chunk's memory and minego does not use it.
	ChunkLight bool
	// ChunkHeightmaps keeps the heightmaps of loaded chunks. Off by default.
	ChunkHeightmaps bool
	// NoTeleport stops the built-in Player from answering ClientboundPlayerPosition
	// (updating the entity, sending AcceptTeleportation and the follow-up move), for
	// callers that confirm teleports themselves and need control over what is sent
	// right after a server correction.
	NoTeleport bool
}

type ResourcePackPolicy string

const (
	ResourcePackAccept  ResourcePackPolicy = "accept"
	ResourcePackDecline ResourcePackPolicy = "decline"
)

type ProxyConfig struct {
	Type     string `json:"type" toml:"type"`
	Host     string `json:"host" toml:"host"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
}

type ConnectOptions struct {
	FakeHost string       `json:"fake_host,omitempty" toml:"fake_host,omitempty"`
	Proxy    *ProxyConfig `json:"proxy,omitempty" toml:"proxy,omitempty"`
	// Keeper, when set, connects through a local connection keeper (see
	// package keeper) instead of dialing the server: the keeper owns the
	// server connection and this client attaches to it, so the client process
	// can restart without the player leaving the server.
	Keeper *KeeperOptions `json:"keeper,omitempty" toml:"keeper,omitempty"`
}

// KeeperOptions locates the local connection keeper and names the session.
type KeeperOptions struct {
	// Socket is the keeper's local socket path.
	Socket string `json:"socket" toml:"socket"`
	// Token is the shared secret the keeper wrote next to its socket.
	Token string `json:"-" toml:"-"`
	// Session names the kept connection, usually the player name. Attaching
	// with the same name resumes the connection if it is still alive.
	Session string `json:"session" toml:"session"`
}

// Detacher is implemented by clients that can leave a kept connection running
// while closing their side (see ConnectOptions.Keeper).
type Detacher interface {
	Detach(ctx context.Context) error
}
