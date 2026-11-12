# Changelog

## v2.0.0 (unreleased)

v2 replaces the packet-handler and string-event APIs with one typed bus, and
moves the module to `github.com/KonjacBot/minego/v2`. It needs Go 1.27, which
added generic methods.

The v0 API stays available under the old import path
`github.com/KonjacBot/minego`. The last v0 commit is `4c5627e`.

### Added

- `pkg/bus`: one bus for packets and events. Its methods are `On`, `OnCtx`,
  `Once`, `WaitFor`, `Stream`, `Emit`, `OnFrame` and `OnRaw`.
  - The handler's parameter type picks what it subscribes to: a clientbound
    packet type such as `*client.KeepAlive`, or any other type as an event.
  - Every subscription returns `off func()`. `OnCtx` removes the
    subscription when its context ends.
- `client.NewClient` now returns `*client.Client`. The client forwards the
  bus methods, so you can call `c.On(...)`, `c.WaitFor(...)` and the others
  directly on it.
- `bot.Client.Bus()`, for code that only holds the interface. Go interfaces
  cannot have generic methods, so such code writes `c.Bus().On(...)`.
- `client.ConnectionStateChangeEvent` can be subscribed to like any other
  event.
- `bot.ClientOptions.ChunkLight` / `ChunkHeightmaps`, and the matching
  `world.WithLight()` / `world.WithHeightmaps()`, keep chunk light and
  heightmaps. Both are dropped by default (see below).

### Changed

- **Handlers run synchronously, in the order they were added.** This covers
  packet, event, raw and frame handlers, and they all run on the read loop.
  - Raw handlers no longer run on goroutines of their own.
  - Long work belongs on your own goroutine.
  - A dispatch slower than 50 ms is logged (`bus.Bus.SlowHandler`).
- **A packet is decoded only if a typed handler subscribed to it.**
  - Packets nobody listens to (sounds, particles, ...) are skipped.
  - A packet that fails to decode only drops the connection when someone
    subscribed to that packet type.
- **Event handlers no longer return an error.** One failing subscriber can no
  longer end the connection. A panicking handler is logged, and the rest still
  run.
- **`OnRaw` payloads are only valid until the handler returns.** The read loop
  reuses its buffer, so copy `p.Data` if you need it later.
- **A decoder panic becomes a connection error** instead of crashing the
  process.
- `Player.CheckServer` waits for the server's next packet by counting packets,
  instead of taking a lock and reading the clock on every packet.
- **The idle read deadline is refreshed only when the read buffer is empty.**
  Before, it was set and cleared once per packet.
- Depends on the go-mc read-path changes (KonjacBot/go-mc#2):
  - Buffered, decrypting reads.
  - Pooled zlib.
  - Decoders that reject length prefixes the packet cannot hold. These used to
    allow allocating 2 GB from 5 bytes, or panicking on negative lengths.

- **World keeps less of each chunk by default.** This cuts a loaded
  overworld column from about 78 KiB to about 26 KiB (-67%). At view distance
  16 (1,089 columns), one bot needs about 28 MB instead of 85 MB.
  - Sky and block light are dropped. They were about half of each chunk,
    nothing in minego read them, and light updates were never applied, so
    they went stale.
  - Heightmaps are dropped too.
  - The spare sections the decoder allocates for the tallest possible column
    are released.
  - Sections that hold a single value (all air, all stone, one biome) share
    one read-only container across all columns and bots. World copies such a
    section before writing to it.
  - World stores its own shallow copy of the column. Other subscribers to
    `*client.LevelChunkWithLight` still see the full packet, light included.
  - `World.SetBlock` and block updates now write through the section instead
    of a copy, so the section's block count is no longer lost.

### Removed

- `bot.AddHandler`, `bot.HandlerFunc` and the `bot.PacketHandler` interface
  (`AddPacketHandler`, `AddRawPacketHandler`, `AddGenericPacketHandler`,
  `HandlePacket`).
- `bot.SubscribeEvent`, `bot.PublishEvent`, `bot.Event` and the
  `bot.EventHandler` interface.
- `bot.Client.PacketHandler()` and `bot.Client.EventHandler()`.
- `client.NewEventHandler` and `client.EventConnectionStateChange`.
- The `EventID() string` methods on event types.

### Performance

`BenchmarkHandleGame` replays 40 ticks of busy play traffic through the real
read loop, with the built-in modules attached (`bench/`). The v0 run already
uses the optimized go-mc:

| | v0 | v2 |
| --- | --- | --- |
| time per packet | 442 ns | 179 ns |
| bytes per run | 105.8 KiB | 26.3 KiB |
| allocations per run | 6133 | 1772 |

Dispatching a packet to two handlers takes 88 ns and makes 0 allocations.

Memory kept per loaded overworld column (`TestLoadedColumnMemory` in
pkg/game/world, generated terrain with vanilla-style light):

| | v0 | v2 default | v2 with light and heightmaps |
| --- | --- | --- | --- |
| bytes per column | 78.2 KiB | 25.9 KiB | 68.6 KiB |
| one bot at view distance 16 (1,089 columns) | about 85 MB | about 28 MB | about 75 MB |

### Migrating from v0

1. Change every import of `github.com/KonjacBot/minego/...` to
   `github.com/KonjacBot/minego/v2/...`. In `go.mod`, require
   `github.com/KonjacBot/minego/v2` and `go 1.27` or later.
2. Keep the `*client.Client` that `client.NewClient` returns. A
   `var c bot.Client` cannot call the generic methods; use `*client.Client`,
   or call `c.Bus().On(...)` instead.
3. Rewrite subscriptions using the table below.
4. Check anything that relied on the old semantics:
   - Raw handlers that expected to run in parallel, or that kept `p.Data`.
   - Event handlers that returned an error to stop the connection.
   - Handlers that block for a long time.

| v0 | v2 |
| --- | --- |
| `bot.AddHandler(c, func(ctx context.Context, p *cp.X) {...})` | `c.On(func(ctx context.Context, p *cp.X) {...})` |
| `bot.SubscribeEvent(c, func(e E) error {...; return nil})` | `c.On(func(ctx context.Context, e E) {...})` |
| `bot.PublishEvent(c, e)` | `c.Emit(ctx, e)` |
| `c.PacketHandler().AddRawPacketHandler(id, f)` | `c.OnRaw(id, f)`, and copy `p.Data` to keep it |
| `c.PacketHandler().AddGenericPacketHandler(f)` | `c.OnFrame(func(id, size) {...})` to notice traffic; subscribe to the specific packet types you need to read |
| `c.PacketHandler().HandlePacket(ctx, p)` (tests) | `c.Bus().Emit(ctx, p)` or `c.Bus().DispatchPacket(ctx, p)` |
| `c.EventHandler().SubscribeEvent(name, f)` | `c.On(func(ctx context.Context, e E) {...})` |
| `func (E) EventID() string` | delete it; the Go type is the key |
| `client.EventConnectionStateChange` | `c.On(func(ctx context.Context, e client.ConnectionStateChangeEvent) {...})` |
| Reading `chunk.Light` / heightmaps from World's columns | set `bot.ClientOptions{ChunkLight: true, ChunkHeightmaps: true}`, or read them from your own `*client.LevelChunkWithLight` handler |
| Polling a flag set by a handler | `v, err := c.WaitFor(ctx, func(v T) bool {...})` |
| A handler registered for one task | `c.OnCtx(taskCtx, ...)`, or call the returned `off()` |
| Test fakes implementing `PacketHandler()` and `EventHandler()` | implement `Bus() *bus.Bus { return b }`, with `b := bus.New()`, and feed packets with `b.Emit(ctx, p)` |

Before and after:

```go
// v0
c := client.NewClient(opts)
bot.SubscribeEvent(c, func(e player.MessageEvent) error {
	fmt.Println(e.Message.String())
	return nil
})
bot.AddHandler(c, func(ctx context.Context, p *cp.SetHealth) {
	if p.Health <= 0 {
		_ = c.WritePacket(ctx, &server.ClientCommand{Action: 0})
	}
})

// v2
c := client.NewClient(opts)
c.On(func(_ context.Context, e player.MessageEvent) {
	fmt.Println(e.Message.String())
})
c.On(func(ctx context.Context, p *cp.SetHealth) {
	if p.Health <= 0 {
		_ = c.WritePacket(ctx, &server.ClientCommand{Action: 0})
	}
})

// New in v2: wait instead of polling, scoped to a task
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()
opened, err := c.WaitFor(ctx, func(e inventory.ContainerOpenEvent) bool { return e.WindowID > 0 })
```

### Migration prompt for AI coding agents

Paste the following into Claude Code or another coding agent at the root of a
project that uses minego v0:

````text
Migrate this Go project from github.com/KonjacBot/minego (v0) to
github.com/KonjacBot/minego/v2. The v2 changelog is at
https://github.com/KonjacBot/minego/blob/main/CHANGELOG.md; follow its
"Migrating from v0" section. Work in this order and keep the build green
between steps:

1. go.mod: require github.com/KonjacBot/minego/v2 (latest) and set
   `go 1.27` or later. Drop the old github.com/KonjacBot/minego requirement
   once nothing imports it. Rewrite every import path
   "github.com/KonjacBot/minego/..." to "github.com/KonjacBot/minego/v2/...".
2. Keep the value returned by client.NewClient as *client.Client. Replace
   variables, fields and parameters typed bot.Client with *client.Client
   where they are used to subscribe. Where only the bot.Client interface is
   available (e.g. code that also accepts test fakes), use c.Bus().On(...)
   instead of c.On(...).
3. Replace subscriptions:
   - bot.AddHandler(c, func(ctx context.Context, p *T) {...})
       -> c.On(func(ctx context.Context, p *T) {...})
   - bot.SubscribeEvent(c, func(e E) error {...})
       -> c.On(func(ctx context.Context, e E) {...}); remove the error
          return. If a handler returned a non-nil error to stop the bot,
          cancel the bot's context (or call c.Close) instead, and say so in
          your summary.
   - bot.PublishEvent(c, e) -> c.Emit(ctx, e), using the nearest context.
   - c.PacketHandler().AddRawPacketHandler(id, f) -> c.OnRaw(id, f). Raw
     handlers are now synchronous and p.Data is only valid during the call:
     copy it (bytes.Clone) if it is stored or used from another goroutine.
   - c.PacketHandler().AddGenericPacketHandler(f) -> if f only notices
     traffic, use c.OnFrame(func(id packetid.ClientboundPacketID, size int)
     {...}); if it reads packet contents, subscribe to those packet types
     with c.On instead (a generic handler forced every packet to be decoded).
   - Delete `EventID() string` methods on event types.
   - client.EventConnectionStateChange string subscriptions
       -> c.On(func(ctx context.Context, e client.ConnectionStateChangeEvent) {...})
4. Fix lifetimes: any handler registered inside a task, request or loop
   iteration (not once at startup) must be removed when that work ends. Use
   c.OnCtx(taskCtx, ...) or keep the returned off func and defer it.
   Replace "set a flag in a handler, then poll it" with
   v, err := c.WaitFor(ctx, func(v T) bool {...}).
5. Handlers now run synchronously on the read loop in registration order.
   Move anything slow (sleeping, pathfinding, waiting on another packet,
   disk or network I/O) into a goroutine that observes the context;
   otherwise it stalls the connection and is logged as a slow dispatch.
   Never call c.WaitFor from inside a handler for a packet that the same
   read loop must deliver: it would deadlock.
6. Tests: fakes that implemented PacketHandler()/EventHandler() should
   implement Bus() *bus.Bus returning a bus.New(); feed decoded packets with
   b.Emit(ctx, packet) and observe events with b.On.
7. Run `go build ./...`, `go vet ./...` and `go test -race ./...` and fix
   everything. Then list every place where behaviour may have changed
   (removed error returns, raw handlers that kept data, handlers that
   blocked, subscriptions you scoped with OnCtx) so a human can review them.
````
