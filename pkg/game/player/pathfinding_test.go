package player

import (
	"errors"
	"math"
	"testing"

	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/go-gl/mathgl/mgl64"
)

func TestAStarRejectsInvalidRequests(t *testing.T) {
	if _, err := AStar(nil, mgl64.Vec3{}, mgl64.Vec3{}, 1); !errors.Is(err, ErrInvalidPathRequest) {
		t.Fatalf("AStar(nil) error = %v, want ErrInvalidPathRequest", err)
	}
	if _, err := AStar(pathTestWorld{}, mgl64.Vec3{math.NaN(), 0, 0}, mgl64.Vec3{}, 1); !errors.Is(err, ErrInvalidPathRequest) {
		t.Fatalf("AStar(NaN) error = %v, want ErrInvalidPathRequest", err)
	}
	if _, err := AStar(pathTestWorld{}, mgl64.Vec3{}, mgl64.Vec3{}, 0); !errors.Is(err, ErrInvalidPathRequest) {
		t.Fatalf("AStar(maxNodes=0) error = %v, want ErrInvalidPathRequest", err)
	}
}

func TestPathMathDoesNotWrapExtremeCoordinates(t *testing.T) {
	a := protocol.Position{math.MinInt32, 0, 0}
	b := protocol.Position{math.MaxInt32, 0, 0}
	if got := heuristic(a, b); got != float64(uint64(math.MaxUint32)) {
		t.Fatalf("heuristic(extremes) = %v, want %v", got, float64(uint64(math.MaxUint32)))
	}
	for _, neighbor := range getNeighbors(b) {
		if neighbor[0] == math.MinInt32 {
			t.Fatal("getNeighbors() wrapped MaxInt32 to MinInt32")
		}
	}
}
