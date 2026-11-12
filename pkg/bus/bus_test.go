package bus

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

type testEvent struct{ N int }

func quiet(b *Bus) *Bus {
	b.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	return b
}

func TestOnPacketAndOff(t *testing.T) {
	b := New()
	var got []int64
	off := b.On(func(_ context.Context, p *client.KeepAlive) { got = append(got, p.ID) })
	if !b.WantsPacket(packetid.ClientboundKeepAlive) {
		t.Fatal("WantsPacket = false with a subscriber")
	}
	if b.WantsPacket(packetid.ClientboundLogin) {
		t.Fatal("WantsPacket = true for a packet nobody subscribed to")
	}
	b.DispatchPacket(context.Background(), &client.KeepAlive{ID: 1})
	off()
	off() // second call is a no-op
	b.DispatchPacket(context.Background(), &client.KeepAlive{ID: 2})
	if len(got) != 1 || got[0] != 1 {
		t.Fatalf("got %v, want [1]", got)
	}
	if b.WantsPacket(packetid.ClientboundKeepAlive) || b.Subscribers() != 0 {
		t.Fatalf("after off: wants=%v subscribers=%d", b.WantsPacket(packetid.ClientboundKeepAlive), b.Subscribers())
	}
}

func TestEventsRunInOrderAndKeepTypesApart(t *testing.T) {
	b := New()
	var order []string
	b.On(func(_ context.Context, e testEvent) { order = append(order, "a") })
	b.On(func(_ context.Context, e testEvent) { order = append(order, "b") })
	b.On(func(_ context.Context, e *testEvent) { order = append(order, "pointer") })
	b.Emit(context.Background(), testEvent{N: 1})
	if len(order) != 2 || order[0] != "a" || order[1] != "b" {
		t.Fatalf("order = %v, want [a b]", order)
	}
}

func TestEmitPacketReachesPacketSubscribers(t *testing.T) {
	b := New()
	var got int64
	b.On(func(_ context.Context, p *client.KeepAlive) { got = p.ID })
	b.Emit(context.Background(), &client.KeepAlive{ID: 7})
	if got != 7 {
		t.Fatalf("got %d, want 7", got)
	}
}

func TestPanicIsIsolated(t *testing.T) {
	b := quiet(New())
	ran := false
	b.On(func(context.Context, testEvent) { panic("boom") })
	b.On(func(context.Context, testEvent) { ran = true })
	b.Emit(context.Background(), testEvent{})
	if !ran {
		t.Fatal("handler after a panicking one did not run")
	}
}

func TestOnCtxRemovesWhenDone(t *testing.T) {
	b := New()
	ctx, cancel := context.WithCancel(context.Background())
	b.OnCtx(ctx, func(context.Context, testEvent) {})
	if b.Subscribers() != 1 {
		t.Fatalf("subscribers = %d, want 1", b.Subscribers())
	}
	cancel()
	deadline := time.Now().Add(time.Second)
	for b.Subscribers() != 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if b.Subscribers() != 0 {
		t.Fatalf("subscribers = %d after cancel, want 0", b.Subscribers())
	}
}

func TestOnceFiresOnce(t *testing.T) {
	b := New()
	n := 0
	b.Once(func(context.Context, testEvent) { n++ })
	for range 3 {
		b.Emit(context.Background(), testEvent{})
	}
	if n != 1 || b.Subscribers() != 0 {
		t.Fatalf("fired %d times, %d subscribers left", n, b.Subscribers())
	}
}

func TestWaitFor(t *testing.T) {
	b := New()
	go func() {
		for i := range 5 {
			time.Sleep(time.Millisecond)
			b.Emit(context.Background(), &client.SetHealth{Health: float32(20 - i*5)})
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p, err := b.WaitFor(ctx, func(p *client.SetHealth) bool { return p.Health < 8 })
	if err != nil {
		t.Fatal(err)
	}
	if p.Health >= 8 {
		t.Fatalf("health = %v", p.Health)
	}
	if b.Subscribers() != 0 {
		t.Fatalf("WaitFor left %d subscribers", b.Subscribers())
	}
}

func TestWaitForTimesOut(t *testing.T) {
	b := New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := b.WaitFor[testEvent](ctx, nil); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
	if b.Subscribers() != 0 {
		t.Fatalf("WaitFor left %d subscribers", b.Subscribers())
	}
}

func TestStreamYieldsInOrderAndUnsubscribes(t *testing.T) {
	ctx := context.Background()
	// Subscription happens when the loop starts, so emit only after that.
	b2 := New()
	var seen []int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for e := range b2.Stream[testEvent](ctx) {
			seen = append(seen, e.N)
			if len(seen) == 5 {
				return
			}
		}
	}()
	for b2.Subscribers() == 0 {
		time.Sleep(time.Millisecond)
	}
	for i := range 5 {
		b2.Emit(ctx, testEvent{N: i})
	}
	<-done
	for i, n := range seen {
		if n != i {
			t.Fatalf("seen = %v, want 0..4 in order", seen)
		}
	}
	if b2.Subscribers() != 0 {
		t.Fatalf("stream left %d subscribers", b2.Subscribers())
	}
}

func TestFrameAndRaw(t *testing.T) {
	b := New()
	var frames, raws int
	offF := b.OnFrame(func(id packetid.ClientboundPacketID, size int) { frames += size })
	offR := b.OnRaw(packetid.ClientboundKeepAlive, func(_ context.Context, p pk.Packet) { raws += len(p.Data) })
	b.DispatchFrame(packetid.ClientboundKeepAlive, 8)
	b.DispatchRaw(context.Background(), packetid.ClientboundKeepAlive, pk.Packet{Data: make([]byte, 8)})
	b.DispatchRaw(context.Background(), packetid.ClientboundLogin, pk.Packet{Data: make([]byte, 99)})
	offF()
	offR()
	if frames != 8 || raws != 8 || b.Subscribers() != 0 {
		t.Fatalf("frames=%d raws=%d subscribers=%d", frames, raws, b.Subscribers())
	}
}

func TestDispatchDoesNotAllocate(t *testing.T) {
	b := New()
	n := 0
	b.On(func(_ context.Context, p *client.KeepAlive) { n++ })
	b.On(func(_ context.Context, p *client.KeepAlive) { n++ })
	b.OnFrame(func(packetid.ClientboundPacketID, int) { n++ })
	ctx := context.Background()
	p := &client.KeepAlive{ID: 1}
	allocs := testing.AllocsPerRun(1000, func() {
		b.DispatchFrame(packetid.ClientboundKeepAlive, 8)
		_ = b.WantsPacket(packetid.ClientboundKeepAlive)
		b.DispatchPacket(ctx, p)
	})
	if allocs != 0 {
		t.Fatalf("dispatch allocated %v times per packet", allocs)
	}
}

// TestSubscribeDuringDispatch exercises copy-on-write under -race: handlers
// add and remove subscriptions while other goroutines dispatch.
func TestSubscribeDuringDispatch(t *testing.T) {
	b := New()
	ctx := context.Background()
	b.On(func(context.Context, *client.KeepAlive) {
		off := b.On(func(context.Context, *client.KeepAlive) {})
		off()
	})
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for i := range 500 {
				b.DispatchPacket(ctx, &client.KeepAlive{ID: int64(i)})
				off := b.On(func(context.Context, testEvent) {})
				b.Emit(ctx, testEvent{N: i})
				off()
			}
		})
	}
	wg.Wait()
	if b.Subscribers() != 1 {
		t.Fatalf("subscribers = %d, want 1", b.Subscribers())
	}
}

func BenchmarkDispatchPacket(b *testing.B) {
	bus := New()
	bus.On(func(context.Context, *client.KeepAlive) {})
	bus.On(func(context.Context, *client.KeepAlive) {})
	ctx := context.Background()
	p := &client.KeepAlive{ID: 1}
	b.ReportAllocs()
	for range b.N {
		if bus.WantsPacket(packetid.ClientboundKeepAlive) {
			bus.DispatchPacket(ctx, p)
		}
	}
}
