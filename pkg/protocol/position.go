package protocol

import "math"

type Position [3]int32

func (p Position) DistanceTo(other Position) float64 {
	dx := float64(p[0] - other[0])
	dy := float64(p[1] - other[1])
	dz := float64(p[2] - other[2])
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (p Position) DistanceToSquared(other Position) float64 {
	dx := float64(p[0] - other[0])
	dy := float64(p[1] - other[1])
	dz := float64(p[2] - other[2])
	return dx*dx + dy*dy + dz*dz
}

func (p Position) Add(other Position) Position {
	return Position{p[0] + other[0], p[1] + other[1], p[2] + other[2]}
}

func (p Position) Sub(other Position) Position {
	return Position{p[0] - other[0], p[1] - other[1], p[2] - other[2]}
}

func (p Position) Mul(scalar float64) Position {
	return Position{int32(float64(p[0]) * scalar), int32(float64(p[1]) * scalar), int32(float64(p[2]) * scalar)}
}

func (p Position) Div(scalar float64) Position {
	return Position{int32(float64(p[0]) / scalar), int32(float64(p[1]) / scalar), int32(float64(p[2]) / scalar)}
}

func (p Position) IsZero() bool {
	return p[0] == 0 && p[1] == 0 && p[2] == 0
}

func (p Position) Clone() Position {
	return Position{p[0], p[1], p[2]}
}

func (p Position) String() string {
	return "(" + string(p[0]) + ", " + string(p[1]) + ", " + string(p[2]) + ")"
}

func (p Position) Equals(other Position) bool {
	return p[0] == other[0] && p[1] == other[1] && p[2] == other[2]
}
