package component

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/net/packet"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

type ChickenVariant struct {
	Mode        int8
	VariantName packet.Identifier
	VariantID   int32
}

func (c *ChickenVariant) ReadFrom(r io.Reader) (n int64, err error) {
	n1, err := (*packet.Byte)(&c.Mode).ReadFrom(r)
	if err != nil {
		return n1, err
	}
	if c.Mode == 0 {
		n2, err := c.VariantName.ReadFrom(r)
		return n1 + n2, err
	}

	if c.Mode == 1 {
		n2, err := (*packet.VarInt)(&c.VariantID).ReadFrom(r)
		return n1 + n2, err
	}

	return n1, err
}

func (c ChickenVariant) WriteTo(w io.Writer) (int64, error) {
	n1, err := (*packet.Byte)(&c.Mode).WriteTo(w)
	if err != nil {
		return n1, err
	}

	if c.Mode == 0 {
		n2, err := c.VariantName.WriteTo(w)
		return n1 + n2, err
	}

	if c.Mode == 1 {
		n2, err := (*packet.VarInt)(&c.VariantID).WriteTo(w)
		return n1 + n2, err
	}
	return n1, err
}

func (*ChickenVariant) Type() slot.ComponentID {
	return 86
}

func (*ChickenVariant) ID() string {
	return "minecraft:chicken/variant"
}
