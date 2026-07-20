package component

import (
	"io"

	"github.com/KonjacBot/go-mc/level/item"
	"github.com/KonjacBot/go-mc/net/packet"
)

func writeNBTValue(w io.Writer, value any) (int64, error) {
	return packet.NBT(value).WriteTo(w)
}

func readNBTValue(r io.Reader, value any) (int64, error) {
	return packet.NBT(value).ReadFrom(r)
}

func (c *GlobalPosition) ReadFrom(r io.Reader) (n int64, err error) {
	var temp int64
	temp, err = (&c.Dimension).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (&c.Position).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	return n, nil
}

func (c GlobalPosition) WriteTo(w io.Writer) (n int64, err error) {
	var temp int64
	temp, err = (&c.Dimension).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	temp, err = (&c.Position).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	return n, nil
}

func decodePotDecorationSide(ids []int32, index int) *item.ID {
	if index >= len(ids) {
		return nil
	}
	value := item.ID(ids[index])
	if value == (item.Brick{}).ID() {
		return nil
	}
	return &value
}

func encodePotDecorationSide(value *item.ID) int32 {
	if value == nil {
		return int32((item.Brick{}).ID())
	}
	return int32(*value)
}
