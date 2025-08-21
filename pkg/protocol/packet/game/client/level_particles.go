package client

import "git.konjactw.dev/patyhank/minego/pkg/protocol/particle"

//codec:gen
type Particle struct {
	LongDistance              bool
	AlwaysVisible             bool
	X, Y, Z                   float64
	OffsetX, OffsetY, OffsetZ float32
	MaxSpeed                  float32
	Count                     int32
	Particle                  particle.Particle
}
