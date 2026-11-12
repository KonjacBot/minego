// Package bus is the single place where a bot's packets and events are
// delivered to the code that wants them.
//
// Subscribe with the generic methods; the type parameter is inferred from the
// handler, so there is nothing to spell out:
//
//	off := b.On(func(ctx context.Context, p *client.KeepAlive) { ... })
//	defer off()
//
//	b.On(func(ctx context.Context, e player.MessageEvent) { ... })
//
// A handler type that implements client.ClientboundPacket subscribes to that
// packet; any other type subscribes to events of that type, published with
// Emit.
//
// Semantics:
//   - Handlers run synchronously on the goroutine that dispatches (for packets,
//     the connection's read loop), in the order they were added. Long work
//     belongs on a goroutine of the handler's own; a dispatch whose handlers
//     together take longer than SlowHandler is logged.
//   - A panicking handler is logged and skipped; the others still run.
//   - Every subscription returns a function that removes it. OnCtx removes it
//     when a context ends, so nothing outlives the work that added it.
//   - Packets nobody subscribed to are never decoded (see WantsPacket).
package bus

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"reflect"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

// DefaultSlowHandler is the default for Bus.SlowHandler.
const DefaultSlowHandler = 50 * time.Millisecond

// StreamBuffer is how many values a Stream holds for a busy loop body before
// the dispatching goroutine has to wait.
const StreamBuffer = 64

const packetSlots = int(packetid.ClientboundPacketIDGuard)

type (
	packetFunc = func(context.Context, client.ClientboundPacket)
	rawFunc    = func(context.Context, pk.Packet)
	eventFunc  = func(context.Context, any)
	// FrameFunc observes every packet as it arrives, before any decoding.
	FrameFunc = func(id packetid.ClientboundPacketID, size int)
)

type entry[F any] struct {
	id uint64
	fn F
}

// list is an immutable snapshot of one topic's subscribers. Dispatch loads it
// without locking; changes replace it under Bus.mu (copy on write), so a
// handler that subscribes or unsubscribes while being dispatched is safe.
type list[F any] []entry[F]

func (l *list[F]) with(id uint64, fn F) *list[F] {
	var next list[F]
	if l != nil {
		next = slices.Clip(slices.Clone(*l))
	}
	next = append(next, entry[F]{id: id, fn: fn})
	return &next
}

// without returns l minus subscription id (nil when empty) and whether id was there.
func (l *list[F]) without(id uint64) (*list[F], bool) {
	if l == nil {
		return nil, false
	}
	i := slices.IndexFunc(*l, func(e entry[F]) bool { return e.id == id })
	if i < 0 {
		return l, false
	}
	if len(*l) == 1 {
		return nil, true
	}
	next := slices.Delete(slices.Clone(*l), i, i+1)
	return &next, true
}

type topic[F any] struct{ head atomic.Pointer[list[F]] }

func (t *topic[F]) load() list[F] {
	if l := t.head.Load(); l != nil {
		return *l
	}
	return nil
}

// Bus delivers packets and events to subscribers. Create one with New.
type Bus struct {
	// SlowHandler logs a warning when the handlers of one packet or event
	// together run longer than this. Zero disables the check.
	SlowHandler time.Duration
	// Logger receives handler panics and slow-handler warnings; nil means
	// slog.Default().
	Logger *slog.Logger

	mu     sync.Mutex // serializes subscription changes
	nextID atomic.Uint64
	count  atomic.Int64

	packets [packetSlots]topic[packetFunc]
	raw     [packetSlots]topic[rawFunc]
	frames  topic[FrameFunc]
	events  atomic.Pointer[map[reflect.Type]*topic[eventFunc]]
}

// New returns an empty Bus.
func New() *Bus {
	return &Bus{SlowHandler: DefaultSlowHandler}
}

func add[F any](b *Bus, t *topic[F], id uint64, fn F) (off func()) {
	b.mu.Lock()
	t.head.Store(t.head.Load().with(id, fn))
	b.mu.Unlock()
	b.count.Add(1)
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if next, removed := t.head.Load().without(id); removed {
			t.head.Store(next)
			b.count.Add(-1)
		}
	}
}

// eventTopic returns the topic for event type key, creating it if needed.
// Topics are never removed, so a topic pointer stays valid once handed out.
func (b *Bus) eventTopic(key reflect.Type) *topic[eventFunc] {
	if m := b.events.Load(); m != nil {
		if t := (*m)[key]; t != nil {
			return t
		}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	old := b.events.Load()
	if old != nil {
		if t := (*old)[key]; t != nil {
			return t
		}
	}
	next := map[reflect.Type]*topic[eventFunc]{}
	if old != nil {
		next = maps.Clone(*old)
	}
	t := new(topic[eventFunc])
	next[key] = t
	b.events.Store(&next)
	return t
}

// packetIDOf reports the packet ID T subscribes to, if T is a packet type.
func packetIDOf[T any]() (packetid.ClientboundPacketID, bool) {
	var zero T
	p, ok := any(zero).(client.ClientboundPacket)
	if !ok {
		return 0, false
	}
	id := p.PacketID()
	if id < 0 || int(id) >= packetSlots {
		panic(fmt.Sprintf("bus: %T has packet ID %d outside the known range", zero, id))
	}
	return id, true
}

// On subscribes f to packets or events of type T and returns a function that
// removes the subscription. Calling the returned function again is harmless.
func (b *Bus) On[T any](f func(ctx context.Context, v T)) (off func()) {
	return b.subscribe[T](b.nextID.Add(1), f)
}

func (b *Bus) subscribe[T any](id uint64, f func(ctx context.Context, v T)) (off func()) {
	if pid, ok := packetIDOf[T](); ok {
		return add(b, &b.packets[pid], id, func(ctx context.Context, p client.ClientboundPacket) {
			f(ctx, p.(T))
		})
	}
	return add(b, b.eventTopic(reflect.TypeFor[T]()), id, func(ctx context.Context, v any) {
		f(ctx, v.(T))
	})
}

// OnCtx is On, removed automatically when ctx is done.
func (b *Bus) OnCtx[T any](ctx context.Context, f func(ctx context.Context, v T)) (off func()) {
	remove := b.On[T](f)
	stop := context.AfterFunc(ctx, remove)
	return func() {
		stop()
		remove()
	}
}

// Once is On for a single delivery: the subscription is removed before f runs.
func (b *Bus) Once[T any](f func(ctx context.Context, v T)) (off func()) {
	var (
		fired  atomic.Bool
		remove func()
		ready  = make(chan struct{})
	)
	remove = b.subscribe[T](b.nextID.Add(1), func(ctx context.Context, v T) {
		if !fired.CompareAndSwap(false, true) {
			return
		}
		<-ready // remove is assigned as soon as subscribe returns
		remove()
		f(ctx, v)
	})
	close(ready)
	return remove
}

// WaitFor blocks until a packet or event of type T arrives for which match
// returns true (the first one when match is nil), or until ctx is done.
func (b *Bus) WaitFor[T any](ctx context.Context, match func(v T) bool) (T, error) {
	found := make(chan T, 1)
	off := b.On[T](func(_ context.Context, v T) {
		if match == nil || match(v) {
			select {
			case found <- v:
			default: // already have one
			}
		}
	})
	defer off()
	select {
	case v := <-found:
		return v, nil
	case <-ctx.Done():
		var zero T
		return zero, context.Cause(ctx)
	}
}

// Stream yields packets or events of type T in order until ctx is done or the
// loop stops; it subscribes when the loop starts and unsubscribes when it
// ends. Up to StreamBuffer values wait for a busy loop body; beyond that the
// dispatcher waits too, so a stalled consumer slows the connection instead of
// dropping values or growing memory without bound.
func (b *Bus) Stream[T any](ctx context.Context) iter.Seq[T] {
	return func(yield func(T) bool) {
		values := make(chan T, StreamBuffer)
		done := make(chan struct{})
		defer close(done)
		off := b.On[T](func(_ context.Context, v T) {
			select {
			case values <- v:
			case <-done:
			case <-ctx.Done():
			}
		})
		defer off()
		for {
			select {
			case v := <-values:
				if !yield(v) {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

// Emit delivers v to the subscribers of its type. A packet goes to the packet
// subscribers, which is also how tests feed decoded packets to game state.
func (b *Bus) Emit[T any](ctx context.Context, v T) {
	if _, isPacket := packetIDOf[T](); isPacket {
		b.DispatchPacket(ctx, any(v).(client.ClientboundPacket))
		return
	}
	m := b.events.Load()
	if m == nil {
		return
	}
	key := reflect.TypeFor[T]()
	t := (*m)[key]
	if t == nil {
		return
	}
	b.emit(ctx, t.load(), v, key)
}

// OnFrame subscribes f to every incoming packet before it is decoded: the
// cheap way to notice traffic without paying for decoding.
func (b *Bus) OnFrame(f FrameFunc) (off func()) {
	return add(b, &b.frames, b.nextID.Add(1), f)
}

// OnRaw subscribes f to the undecoded payload of packets with the given ID.
// p.Data is only valid until f returns; copy it to keep it.
func (b *Bus) OnRaw(id packetid.ClientboundPacketID, f func(ctx context.Context, p pk.Packet)) (off func()) {
	if id < 0 || int(id) >= packetSlots {
		panic(fmt.Sprintf("bus: packet ID %d outside the known range", id))
	}
	return add(b, &b.raw[id], b.nextID.Add(1), f)
}

// Subscribers reports the number of live subscriptions, for leak checks.
func (b *Bus) Subscribers() int {
	return int(b.count.Load())
}

// WantsPacket reports whether a typed handler subscribes to id, i.e. whether
// the packet needs decoding at all.
func (b *Bus) WantsPacket(id packetid.ClientboundPacketID) bool {
	return id >= 0 && int(id) < packetSlots && b.packets[id].head.Load() != nil
}

// DispatchFrame runs the OnFrame subscribers. They are expected to be
// trivial, so unlike the other dispatchers they are not timed.
func (b *Bus) DispatchFrame(id packetid.ClientboundPacketID, size int) {
	for _, e := range b.frames.load() {
		b.callFrame(e.fn, id, size)
	}
}

// DispatchRaw runs the OnRaw subscribers of id.
func (b *Bus) DispatchRaw(ctx context.Context, id packetid.ClientboundPacketID, p pk.Packet) {
	if id < 0 || int(id) >= packetSlots {
		return
	}
	l := b.raw[id].load()
	if l == nil {
		return
	}
	start := b.start()
	for _, e := range l {
		b.callRaw(e.fn, ctx, p, id)
	}
	b.checkSlow(start, id)
}

// DispatchPacket runs the typed subscribers of a decoded packet.
func (b *Bus) DispatchPacket(ctx context.Context, p client.ClientboundPacket) {
	id := p.PacketID()
	if id < 0 || int(id) >= packetSlots {
		return
	}
	l := b.packets[id].load()
	if l == nil {
		return
	}
	start := b.start()
	for _, e := range l {
		b.callPacket(e.fn, ctx, p, id)
	}
	b.checkSlow(start, id)
}

func (b *Bus) emit(ctx context.Context, l list[eventFunc], v any, key reflect.Type) {
	start := b.start()
	for _, e := range l {
		b.callEvent(e.fn, ctx, v, key)
	}
	b.checkSlow(start, key)
}

func (b *Bus) callPacket(f packetFunc, ctx context.Context, p client.ClientboundPacket, id packetid.ClientboundPacketID) {
	defer b.recoverHandler(id)
	f(ctx, p)
}

func (b *Bus) callRaw(f rawFunc, ctx context.Context, p pk.Packet, id packetid.ClientboundPacketID) {
	defer b.recoverHandler(id)
	f(ctx, p)
}

func (b *Bus) callFrame(f FrameFunc, id packetid.ClientboundPacketID, size int) {
	defer b.recoverHandler(id)
	f(id, size)
}

func (b *Bus) callEvent(f eventFunc, ctx context.Context, v any, key reflect.Type) {
	defer b.recoverHandler(key)
	f(ctx, v)
}

// recoverHandler is deferred around every handler and turns a panic into a
// log entry, so the remaining handlers still run.
func (b *Bus) recoverHandler(what any) {
	if r := recover(); r != nil {
		b.logger().Error("bus: handler panicked", "for", what, "panic", r, "stack", string(debug.Stack()))
	}
}

// start and checkSlow time one dispatch (all handlers of one packet or event)
// rather than each handler, which keeps clock reads to two per dispatch.
func (b *Bus) start() time.Time {
	if b.SlowHandler > 0 {
		return time.Now()
	}
	return time.Time{}
}

func (b *Bus) checkSlow(start time.Time, what any) {
	if start.IsZero() {
		return
	}
	if d := time.Since(start); d > b.SlowHandler {
		b.logger().Warn("bus: slow handlers held up dispatch", "for", what, "took", d)
	}
}

func (b *Bus) logger() *slog.Logger {
	if b.Logger != nil {
		return b.Logger
	}
	return slog.Default()
}
