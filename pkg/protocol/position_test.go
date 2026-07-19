package protocol

import (
	"math"
	"testing"
)

func TestPositionDistanceDoesNotOverflowInt32(t *testing.T) {
	a := Position{math.MaxInt32, 0, 0}
	b := Position{math.MinInt32, 0, 0}
	want := float64(math.MaxUint32)
	if got := a.DistanceTo(b); got != want {
		t.Fatalf("DistanceTo() = %v, want %v", got, want)
	}
	if got := a.DistanceToSquared(b); got != want*want {
		t.Fatalf("DistanceToSquared() = %v, want %v", got, want*want)
	}
}
