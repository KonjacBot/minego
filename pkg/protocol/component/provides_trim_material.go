package component

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type ProvidesTrimMaterial struct {
	Mode     int8
	Name     packet.Identifier
	Material packet.OptID[TrimMaterial, *TrimMaterial]
}

func (p *ProvidesTrimMaterial) ReadFrom(r io.Reader) (n int64, err error) {
	n1, err := (*packet.Byte)(&p.Mode).ReadFrom(r)
	if err != nil {
		return n1, err
	}
	if p.Mode == 0 {
		n2, err := p.Name.ReadFrom(r)
		return n1 + n2, err
	}

	if p.Mode == 1 {
		n2, err := p.Material.ReadFrom(r)
		return n1 + n2, err
	}

	return n1, err
}

func (p ProvidesTrimMaterial) WriteTo(w io.Writer) (int64, error) {
	n1, err := (*packet.Byte)(&p.Mode).WriteTo(w)
	if err != nil {
		return n1, err
	}

	if p.Mode == 0 {
		n2, err := p.Name.WriteTo(w)
		return n1 + n2, err
	}

	if p.Mode == 1 {
		n2, err := p.Material.WriteTo(w)
		return n1 + n2, err
	}
	return n1, err
}

func (*ProvidesTrimMaterial) Type() slot.ComponentID {
	return 53
}

func (*ProvidesTrimMaterial) ID() string {
	return "minecraft:provides_trim_material"
}
