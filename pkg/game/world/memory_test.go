package world

import (
	"bytes"
	"context"
	"math/rand/v2"
	"runtime"
	"testing"

	"github.com/KonjacBot/go-mc/level"
	"github.com/KonjacBot/go-mc/level/block"

	"github.com/KonjacBot/minego/v2/pkg/protocol"
	gameclient "github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
)

// Overworld columns have 24 sections (y -64..319).
const overworldSections = 24

// typicalColumn builds a chunk column shaped like generated overworld terrain:
// deepslate with ores and lava at the bottom, stone with ores and caves above
// it, a surface section with grass, water and trees, a section of tree tops,
// and air above. Light is what a vanilla server sends: full sky light over the
// surface, none underground, block light around the lava.
func typicalColumn(seed uint64) *level.Chunk {
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	c := level.EmptyChunk(overworldSections)
	pick := func(states []level.BlocksState) level.BlocksState { return states[rng.IntN(len(states))] }

	deep := []level.BlocksState{1, 1, 1, 1, 1, 1, 2, 3, 4, 5, 6, 7, 0}     // mostly one block, some ores, caves
	stone := []level.BlocksState{8, 8, 8, 8, 8, 9, 10, 11, 12, 13, 14, 15, 0} // stone, granite, ores, gravel, air
	surface := make([]level.BlocksState, 0, 28)
	for s := range level.BlocksState(28) {
		surface = append(surface, 20+s) // grass, dirt, flowers, logs, leaves in several states, water levels
	}
	treeTops := []level.BlocksState{0, 0, 0, 0, 0, 0, 0, 30, 31, 32, 33}

	full := bytes.Repeat([]byte{0xff}, 2048)
	for i := range c.Sections {
		s := &c.Sections[i]
		var states []level.BlocksState
		switch {
		case i < 4:
			states = deep
			s.BlockLight = make([]byte, 2048) // lava light
			for j := range 300 {
				s.BlockLight[j*5] = 0xdc
			}
		case i < 8:
			states = stone
		case i == 8:
			states = surface
			s.SkyLight = make([]byte, 2048)
			for j := range s.SkyLight {
				s.SkyLight[j] = byte(rng.IntN(256))
			}
		case i == 9:
			states = treeTops
			s.SkyLight = bytes.Clone(full)
		default:
			s.SkyLight = bytes.Clone(full)
			continue // air
		}
		for b := range 4096 {
			s.SetBlock(b, pick(states))
		}
	}
	for range 6 {
		c.BlockEntity = append(c.BlockEntity, level.BlockEntity{XZ: int8(rng.IntN(256)), Y: int16(rng.IntN(64)), Type: 7})
	}
	return c
}

func columnPacket(tb testing.TB, x, z int32, seed uint64) []byte {
	tb.Helper()
	var buf bytes.Buffer
	p := gameclient.LevelChunkWithLight{Pos: level.ChunkPos{x, z}, Data: typicalColumn(seed)}
	if _, err := p.WriteTo(&buf); err != nil {
		tb.Fatal(err)
	}
	return buf.Bytes()
}

func heapInUse() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// loadedColumnBytes reports the heap one loaded column keeps alive in World,
// averaged over 9x9 columns.
func loadedColumnBytes(tb testing.TB, opts ...Option) (perColumn float64, wire int) {
	tb.Helper()
	const side = 9
	packets := make([][]byte, 0, side*side)
	for x := range int32(side) {
		for z := range int32(side) {
			p := columnPacket(tb, x, z, uint64(x*side+z))
			wire += len(p)
			packets = append(packets, p)
		}
	}

	c := newWorldTestClient()
	w := NewWorld(c, opts...)
	c.world = w
	ctx := context.Background()

	before := heapInUse()
	for _, raw := range packets {
		var p gameclient.LevelChunkWithLight
		if _, err := p.ReadFrom(bytes.NewReader(raw)); err != nil {
			tb.Fatal(err)
		}
		c.handler.Emit(ctx, &p)
	}
	after := heapInUse()
	runtime.KeepAlive(w)
	runtime.KeepAlive(packets)
	return float64(after-before) / float64(len(packets)), wire / len(packets)
}

// maxColumnBytes guards the default memory footprint of a loaded column.
// v0 kept about 78 KiB per column of this terrain; v2 keeps about 26 KiB.
const maxColumnBytes = 35 << 10

func TestLoadedColumnMemory(t *testing.T) {
	perColumn, wire := loadedColumnBytes(t)
	t.Logf("one loaded overworld column keeps %.0f bytes (%.1f KiB) alive; it arrived as %d bytes",
		perColumn, perColumn/1024, wire)
	if perColumn > maxColumnBytes {
		t.Fatalf("a loaded column keeps %.0f bytes, over the %d byte budget", perColumn, maxColumnBytes)
	}
}

func TestLoadedColumnMemoryWithLight(t *testing.T) {
	perColumn, _ := loadedColumnBytes(t, WithLight(), WithHeightmaps())
	t.Logf("with light and heightmaps kept: %.0f bytes (%.1f KiB) per column", perColumn, perColumn/1024)
}

// TestSharedSectionsAreCopiedOnWrite loads two columns whose upper sections
// are all air, so they share one container, then writes a block into one of
// them. The other column, and the shared container, must stay air.
func TestSharedSectionsAreCopiedOnWrite(t *testing.T) {
	c := newWorldTestClient()
	w := NewWorld(c)
	c.world = w
	ctx := context.Background()
	for _, x := range []int32{0, 1} {
		var p gameclient.LevelChunkWithLight
		if _, err := p.ReadFrom(bytes.NewReader(columnPacket(t, x, 0, uint64(x)))); err != nil {
			t.Fatal(err)
		}
		c.handler.Emit(ctx, &p)
	}

	a := w.columns[level.ChunkPos{0, 0}]
	b := w.columns[level.ChunkPos{1, 0}]
	const airSection = 20 // y 256..271
	if a.Sections[airSection].States != b.Sections[airSection].States {
		t.Fatal("all-air sections of two columns do not share a container")
	}
	shared := a.Sections[airSection].States

	pos := protocol.Position{3, 260, 5}
	stone := block.StateList[1]
	if err := w.SetBlock(pos, stone); err != nil {
		t.Fatal(err)
	}
	if got, _ := w.GetBlock(pos); got != stone {
		t.Fatalf("written block = %v, want %v", got, stone)
	}
	if got, _ := w.GetBlock(protocol.Position{16 + 3, 260, 5}); got.ID() != "minecraft:air" {
		t.Fatalf("other column's block = %v, want air", got)
	}
	if v, ok := shared.SingleValue(); !ok || v != 0 {
		t.Fatal("the shared container was modified")
	}
}
