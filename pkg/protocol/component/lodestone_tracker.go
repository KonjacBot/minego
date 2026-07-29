package component

import (
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type LodestoneTracker struct {
	HasGlobalPosition bool
	Dimension         pk.Option[pk.Identifier, *pk.Identifier]
	Position          pk.Option[pk.Position, *pk.Position]
	Tracked           bool
}

func (*LodestoneTracker) ID() string {
	return "minecraft:lodestone_tracker"
}

func (c *LodestoneTracker) ReadFrom(r io.Reader) (n int64, err error) {
	*c = LodestoneTracker{}

	var hasGlobalPosition pk.Boolean
	temp, err := hasGlobalPosition.ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	c.HasGlobalPosition = bool(hasGlobalPosition)

	if c.HasGlobalPosition {
		c.Dimension.Has = true
		temp, err = (&c.Dimension.Val).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}

		c.Position.Has = true
		temp, err = (&c.Position.Val).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	}

	temp, err = (*pk.Boolean)(&c.Tracked).ReadFrom(r)
	n += temp
	return n, err
}

func (c LodestoneTracker) WriteTo(w io.Writer) (n int64, err error) {
	temp, err := pk.Boolean(c.HasGlobalPosition).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}

	if c.HasGlobalPosition {
		temp, err = c.Dimension.Val.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}

		temp, err = c.Position.Val.WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}

	temp, err = pk.Boolean(c.Tracked).WriteTo(w)
	n += temp
	return n, err
}
