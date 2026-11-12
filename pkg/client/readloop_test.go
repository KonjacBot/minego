package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"

	mcnet "github.com/KonjacBot/go-mc/net"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/game/player"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

const replayThreshold = 256

// replayConn plays back a recorded server stream and discards writes.
type replayConn struct {
	net.Conn
	r *bytes.Reader
}

func (c *replayConn) Read(p []byte) (int, error)       { return c.r.Read(p) }
func (c *replayConn) Write(p []byte) (int, error)      { return len(p), nil }
func (c *replayConn) Close() error                     { return nil }
func (c *replayConn) SetDeadline(time.Time) error      { return nil }
func (c *replayConn) SetReadDeadline(time.Time) error  { return nil }
func (c *replayConn) SetWriteDeadline(time.Time) error { return nil }

// busyServerTick is one tick of play traffic around a few moving mobs: mostly
// entity movement, plus sounds, time and a chat line.
func busyServerTick(tick int) []client.ClientboundPacket {
	var out []client.ClientboundPacket
	for e := range int32(10) {
		out = append(out, &client.UpdateEntityPosition{EntityID: e, DeltaX: int16(tick), DeltaZ: 64})
	}
	for e := range int32(4) {
		out = append(out,
			&client.UpdateEntityRotation{EntityID: e},
			&client.SetHeadRotation{EntityID: e},
		)
	}
	out = append(out,
		&client.SetEntityVelocity{EntityID: 1},
		&client.SetEntityVelocity{EntityID: 2},
		&client.SoundEffect{SoundCategory: 1, Volume: 1, Pitch: 1},
		&client.SoundEffect{SoundCategory: 2, Volume: 1, Pitch: 1},
		&client.SetTime{GameTime: int64(tick)},
	)
	return out
}

func replayStream(tb testing.TB, ticks int) (stream []byte, packets int) {
	tb.Helper()
	var buf bytes.Buffer
	for tick := range ticks {
		for _, p := range busyServerTick(tick) {
			raw := pk.Marshal(p.PacketID(), p.(pk.FieldEncoder))
			if err := raw.Pack(&buf, replayThreshold); err != nil {
				tb.Fatal(err)
			}
			packets++
		}
	}
	return buf.Bytes(), packets
}

// playOnce runs HandleGame over stream until the stream ends.
func playOnce(tb testing.TB, c *Client, stream []byte) {
	tb.Helper()
	conn := mcnet.WrapConn(&replayConn{r: bytes.NewReader(stream)})
	conn.SetThreshold(replayThreshold)
	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()
	c.setConnectionState(packet.StatePlay)
	if err := c.HandleGame(context.Background()); !errors.Is(err, io.EOF) {
		tb.Fatalf("HandleGame = %v, want EOF at the end of the stream", err)
	}
}

func BenchmarkHandleGame(b *testing.B) {
	stream, packets := replayStream(b, 40)
	b.ReportAllocs()
	b.SetBytes(int64(len(stream)))
	b.ResetTimer()
	for range b.N {
		b.StopTimer()
		c := NewClient(nil)
		b.StartTimer()
		playOnce(b, c, stream)
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*packets), "ns/packet")
}

// TestReconnectDoesNotLeak plays many sessions on one client, each with a
// user subscription scoped to the session, and checks that goroutines and
// subscriptions return to where they started.
func TestReconnectDoesNotLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	stream, _ := replayStream(t, 5)
	c := NewClient(nil)
	baseline := c.Bus().Subscribers()

	for range 100 {
		session, end := context.WithCancel(context.Background())
		moves := 0
		c.OnCtx(session, func(context.Context, *client.UpdateEntityPosition) { moves++ })
		c.OnCtx(session, func(context.Context, player.MessageEvent) {})
		playOnce(t, c, stream)
		if moves != 50 {
			t.Fatalf("saw %d entity moves, want 50", moves)
		}
		end()
		if err := c.Close(context.Background()); err != nil {
			t.Fatal(err)
		}
	}

	deadline := time.Now().Add(time.Second)
	for c.Bus().Subscribers() != baseline && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := c.Bus().Subscribers(); got != baseline {
		t.Fatalf("subscribers = %d after 100 sessions, want %d", got, baseline)
	}
}

// TestUnsubscribedPacketsAreNotDecoded sends a packet whose payload is not a
// valid SoundEffect. Nobody subscribes to SoundEffect, so it must be skipped,
// not decoded and turned into a connection error.
func TestUnsubscribedPacketsAreNotDecoded(t *testing.T) {
	var buf bytes.Buffer
	garbage := pk.Packet{ID: int32((&client.SoundEffect{}).PacketID()), Data: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}
	if err := garbage.Pack(&buf, replayThreshold); err != nil {
		t.Fatal(err)
	}
	c := NewClient(nil)
	playOnce(t, c, buf.Bytes())

	c.On(func(context.Context, *client.SoundEffect) {})
	conn := mcnet.WrapConn(&replayConn{r: bytes.NewReader(buf.Bytes())})
	conn.SetThreshold(replayThreshold)
	c.conn = conn
	c.setConnectionState(packet.StatePlay)
	err := c.HandleGame(context.Background())
	if err == nil || errors.Is(err, io.EOF) {
		t.Fatalf("with a subscriber, HandleGame = %v, want a decode error", err)
	}
	// The failed packet is dumped to a temp file for debugging; remove it.
	if _, dump, ok := strings.Cut(err.Error(), "TEMPPACKET: "); ok && dump != "" {
		_ = os.Remove(dump)
	}
}
