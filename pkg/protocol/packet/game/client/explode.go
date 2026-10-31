package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/pkg/protocol/component"
	"github.com/KonjacBot/minego/pkg/protocol/particle"
)

var _ ClientboundPacket = (*Explode)(nil)

//codec:gen
type Vec3 struct {
	X, Y, Z float64
}

// Explode is the 26.2 explosion packet: center, radius, number of blocks affected, knockback for the
// receiving player, the explosion particle and sound, and the weighted list of block particles.
// Missing the fields after the knockback made every nearby explosion end the connection (unexpected EOF).
//
//codec:gen
type Explode struct {
	CenterX                 float64
	CenterY                 float64
	CenterZ                 float64
	Radius                  float32
	BlockCount              int32
	PlayerKnockbackVelocity pk.Option[Vec3, *Vec3]
	ExplosionParticle       particle.Particle
	ExplosionSound          pk.OptID[component.SoundEvent, *component.SoundEvent]
	BlockParticles          []ExplosionBlockParticle
}

// ExplosionBlockParticle is one weighted alternative for the particles spawned at exploded blocks.
//
//codec:gen
type ExplosionBlockParticle struct {
	Particle particle.Particle
	Scaling  float32
	Speed    float32
	Weight   int32 `mc:"VarInt"`
}

func (Explode) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundExplode
}
