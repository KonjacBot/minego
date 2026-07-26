package player

import (
	"errors"
	"testing"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/KonjacBot/go-mc/data/entity"
	"github.com/KonjacBot/go-mc/level/block"

	"github.com/KonjacBot/minego/pkg/bot"
	"github.com/KonjacBot/minego/pkg/protocol"
)

func BenchmarkAStarOpenPlane80(b *testing.B) {
	world := openPlanePathWorld{}
	start := mgl64.Vec3{0.5, 1, 0.5}
	goal := mgl64.Vec3{80.5, 1, 80.5}
	b.ResetTimer()
	for range b.N {
		path, err := AStar(world, start, goal, 20_000)
		if err != nil {
			b.Fatal(err)
		}
		if len(path) == 0 {
			b.Fatal("no path")
		}
	}
}

func BenchmarkFastAStarWithinOpenPlane80(b *testing.B) {
	world := openPlanePathWorld{}
	start := mgl64.Vec3{0.5, 1, 0.5}
	goal := mgl64.Vec3{80.5, 1.5, 80.5}
	b.ResetTimer()
	for range b.N {
		path, err := FastAStarWithin(world, start, goal, 4096, 4.5)
		if err != nil {
			b.Fatal(err)
		}
		if len(path) == 0 {
			b.Fatal("no path")
		}
	}
}

func BenchmarkFastFlightAStarWithinOpenAir80(b *testing.B) {
	world := openPlanePathWorld{}
	start := mgl64.Vec3{0.5, 3, 0.5}
	goal := mgl64.Vec3{80.5, 3.5, 80.5}
	b.ResetTimer()
	for range b.N {
		path, err := FastFlightAStarWithin(world, start, goal, 4096, 4.5)
		if err != nil {
			b.Fatal(err)
		}
		if len(path) == 0 {
			b.Fatal("no path")
		}
	}
}

func TestFastAStarWithinStopsAtInteractionRange(t *testing.T) {
	world := openPlanePathWorld{}
	goal := mgl64.Vec3{20.5, 1.5, 0.5}
	path, err := FastAStarWithin(world, mgl64.Vec3{0.5, 1, 0.5}, goal, 4096, 4.5)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) == 0 {
		t.Fatal("no path")
	}
	last := vectorCell(path[len(path)-1])
	if !reachedGoal(last, goal, 4.5) {
		t.Fatalf("last cell %v is outside interaction range of %v", last, goal)
	}
	if last == vectorCell(goal) {
		t.Fatalf("range search walked to exact goal %v", last)
	}
}

func TestFastAStarWithinAllowsBlockedExactGoal(t *testing.T) {
	world := openPlanePathWorld{blocked: map[protocol.Position]bool{
		{10, 1, 0}: true,
	}}
	goal := mgl64.Vec3{10.5, 1.5, 0.5}
	path, err := FastAStarWithin(world, mgl64.Vec3{0.5, 1, 0.5}, goal, 4096, 4.5)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) == 0 || !reachedGoal(vectorCell(path[len(path)-1]), goal, 4.5) {
		t.Fatalf("path did not stop near blocked goal: %v", path)
	}
}

func TestFastFlightAStarWithinTraversesUnsupportedAir(t *testing.T) {
	world := openPlanePathWorld{}
	start := mgl64.Vec3{0.5, 3, 0.5}
	goal := mgl64.Vec3{10.5, 3.5, 0.5}

	groundPath, err := FastAStarWithin(world, start, goal, 4096, 4.5)
	if err != nil {
		t.Fatal(err)
	}
	if len(groundPath) != 0 {
		t.Fatalf("ground path unexpectedly crossed unsupported air: %v", groundPath)
	}

	flightPath, err := FastFlightAStarWithin(world, start, goal, 4096, 4.5)
	if err != nil {
		t.Fatal(err)
	}
	if len(flightPath) == 0 || !reachedGoal(vectorCell(flightPath[len(flightPath)-1]), goal, 4.5) {
		t.Fatalf("flight path did not reach target range: %v", flightPath)
	}
}

func TestMinecraftNeighborsStepUpOneBlock(t *testing.T) {
	world := openPlanePathWorld{raised: map[protocol.Position]bool{
		{1, 1, 0}: true,
	}}
	path, err := AStar(world, mgl64.Vec3{0.5, 1, 0.5}, mgl64.Vec3{1.5, 2, 0.5}, 64)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) != 2 || vectorCell(path[1]) != (protocol.Position{1, 2, 0}) {
		t.Fatalf("step-up path = %v", path)
	}
}

func TestMinecraftNeighborsDoNotCutBlockedDiagonal(t *testing.T) {
	world := openPlanePathWorld{blocked: map[protocol.Position]bool{
		{1, 1, 0}: true,
		{0, 1, 1}: true,
	}}
	walkable := newWalkabilityCache(world).walkable
	for _, neighbor := range minecraftNeighbors(protocol.Position{0, 1, 0}, walkable) {
		if neighbor == (protocol.Position{1, 1, 1}) {
			t.Fatalf("diagonal neighbor crossed blocked corners: %v", neighbor)
		}
	}
}

type openPlanePathWorld struct {
	blocked map[protocol.Position]bool
	raised  map[protocol.Position]bool
}

func (world openPlanePathWorld) GetBlock(pos protocol.Position) (block.Block, error) {
	if world.blocked[pos] {
		return block.Stone{}, nil
	}
	if world.raised[pos] {
		return block.Stone{}, nil
	}
	switch pos[1] {
	case 0:
		return block.Stone{}, nil
	case 1, 2, 3, 4:
		return block.Air{}, nil
	default:
		return nil, errors.New("block not loaded")
	}
}
func (openPlanePathWorld) SetBlock(protocol.Position, block.Block) error { return nil }
func (openPlanePathWorld) GetNearbyBlocks(protocol.Position, int32) ([]block.Block, error) {
	return nil, nil
}
func (openPlanePathWorld) FindNearbyBlock(protocol.Position, int32, block.Block) (protocol.Position, error) {
	return protocol.Position{}, errors.New("not implemented")
}
func (openPlanePathWorld) Entities() []bot.Entity               { return nil }
func (openPlanePathWorld) GetEntity(int32) bot.Entity           { return nil }
func (openPlanePathWorld) GetNearbyEntities(int32) []bot.Entity { return nil }
func (openPlanePathWorld) GetEntitiesByType(entity.ID) []bot.Entity {
	return nil
}
