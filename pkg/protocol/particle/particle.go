package particle

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type ParticleData interface {
	io.ReaderFrom
	io.WriterTo
	ParticleID() int32
}

type Particle struct {
	ID   int32
	Data ParticleData
}

func (p *Particle) ReadFrom(r io.Reader) (int64, error) {
	n, err := (*pk.VarInt)(&p.ID).ReadFrom(r)
	if err != nil {
		return n, err
	}

	switch p.ID {
	case 1: // Block
		data := &Block{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 2: // BlockMarker
		data := &BlockMarker{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 13: // Dust
		data := &Dust{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 14: // DustColorTransition
		data := &DustColorTransition{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 20: // EntityEffect
		data := &EntityEffect{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 28: // FallingDust
		data := &FallingDust{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 35: // TintedLeaves
		data := &TintedLeaves{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 37: // SculkCharge
		data := &SculkCharge{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 46: // Item
		data := &Item{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 47: // Vibration
		data := &Vibration{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 48: // Trail
		data := &Trail{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 102: // Shriek
		data := &Shriek{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 108: // DustPillar
		data := &DustPillar{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	case 112: // BlockCrumble
		data := &BlockCrumble{}
		n2, err := data.ReadFrom(r)
		p.Data = data
		return n + n2, err
	default:
		// BasicParticle - no additional data
		p.Data = nil
		return n, nil
	}
}

func (p Particle) WriteTo(w io.Writer) (int64, error) {
	n, err := (*pk.VarInt)(&p.ID).WriteTo(w)
	if err != nil {
		return n, err
	}

	if p.Data == nil {
		return n, nil
	}

	n2, err := p.Data.WriteTo(w)
	return n + n2, err
}
