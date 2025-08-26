package particle

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

// BasicParticle for Particle Type Other
type BasicParticle struct {
	ID int32
}

func (b *BasicParticle) ReadFrom(r io.Reader) (int64, error) {
	return (*pk.VarInt)(&b.ID).ReadFrom(r)
}

func (b BasicParticle) WriteTo(w io.Writer) (int64, error) {
	return (*pk.VarInt)(&b.ID).WriteTo(w)
}

func (b BasicParticle) ParticleID() int32 {
	return b.ID
}

// Block for Particle Type 1
//
//codec:gen
type Block struct {
	BlockState int32 `mc:"VarInt"`
}

func (b Block) ParticleID() int32 {
	return 1
}

// BlockMarker for Particle Type 2
//
//codec:gen
type BlockMarker struct {
	BlockState int32 `mc:"VarInt"`
}

func (b BlockMarker) ParticleID() int32 {
	return 2
}

// Dust for Particle Type 13
//
//codec:gen
type Dust struct {
	Color int32
	Scale float32
}

func (d Dust) ParticleID() int32 {
	return 13
}

// DustColorTransition for Particle Type 14
//
//codec:gen
type DustColorTransition struct {
	FromColor int32
	ToColor   int32
	Scale     float32
}

func (d DustColorTransition) ParticleID() int32 {
	return 14
}

// EntityEffect for Particle Type 20
//
//codec:gen
type EntityEffect struct {
	Color int32
}

func (e EntityEffect) ParticleID() int32 {
	return 20
}

// FallingDust for Particle Type 28
//
//codec:gen
type FallingDust struct {
	BlockState int32 `mc:"VarInt"`
}

func (f FallingDust) ParticleID() int32 {
	return 28
}

// TintedLeaves for Particle Type 35
//
//codec:gen
type TintedLeaves struct {
	Color int32
}

func (t TintedLeaves) ParticleID() int32 {
	return 35
}

// SculkCharge for Particle Type 37
//
//codec:gen
type SculkCharge struct {
	Roll float32
}

func (s SculkCharge) ParticleID() int32 {
	return 37
}

// Item for Particle Type 46
//
//codec:gen
type Item struct {
	Item slot.Slot
}

func (i Item) ParticleID() int32 {
	return 46
}

// Vibration for Particle Type 47
type Vibration struct {
	Type int32 `mc:"VarInt"`
	//if Type eq 0
	BlockPosition pk.Position
	//elif Type eq 1
	EntityID        int32
	EntityEyeHeight float32
	//end
	Ticks int32 `mc:"VarInt"`
}

func (v Vibration) WriteTo(w io.Writer) (n int64, err error) {
	(*pk.VarInt)(&v.Type).WriteTo(w)
	n += 4
	switch v.Type {
	case 0:
		n1, err := v.BlockPosition.WriteTo(w)
		if err != nil {
			return n + n1, err
		}
	case 1:
		n1, err := (*pk.VarInt)(&v.EntityID).WriteTo(w)
		if err != nil {
			return n + n1, err
		}
		n2, err := (*pk.Float)(&v.EntityEyeHeight).WriteTo(w)
		if err != nil {
			return n + n1 + n2, err
		}
		n += n1 + n2
	}
	n2, err := (*pk.VarInt)(&v.Ticks).WriteTo(w)
	if err != nil {
		return n + n2, err
	}
	return n + n2, err
}

func (v *Vibration) ReadFrom(r io.Reader) (int64, error) {
	n, err := (*pk.VarInt)(&v.Type).ReadFrom(r)
	if err != nil {
		return n, err
	}
	switch v.Type {
	case 0:
		n1, err := v.BlockPosition.ReadFrom(r)
		if err != nil {
			return n + n1, err
		}
		return n + n1, err
	case 1:
		n1, err := (*pk.VarInt)(&v.EntityID).ReadFrom(r)
		if err != nil {
			return n + n1, err
		}
		n2, err := (*pk.Float)(&v.EntityEyeHeight).ReadFrom(r)
		if err != nil {
			return n + n1 + n2, err
		}
		return n + n1 + n2, err
	}
	n2, err := (*pk.VarInt)(&v.Ticks).ReadFrom(r)
	if err != nil {
		return n + n2, err
	}
	return n + n2, err
}

func (v Vibration) ParticleID() int32 {
	return 47
}

// Trail for Particle Type 48
//
//codec:gen
type Trail struct {
	X, Y, Z  float64
	Color    int32
	Duration int32 `mc:"VarInt"`
}

func (t Trail) ParticleID() int32 {
	return 48
}

// Shriek for Particle Type 102
//
//codec:gen
type Shriek struct {
	Delay int32 `mc:"VarInt"`
}

func (s Shriek) ParticleID() int32 {
	return 102
}

// DustPillar for Particle Type 108
//
//codec:gen
type DustPillar struct {
	BlockState int32 `mc:"VarInt"`
}

func (d DustPillar) ParticleID() int32 {
	return 108
}

// BlockCrumble for Particle Type 112
//
//codec:gen
type BlockCrumble struct {
	BlockState int32 `mc:"VarInt"`
}

func (b BlockCrumble) ParticleID() int32 {
	return 112
}
