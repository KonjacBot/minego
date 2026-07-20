package client

import (
	"io"

	"github.com/KonjacBot/go-mc/chat"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type MapIcon struct {
	Type        int32 `mc:"VarInt"`
	X, Z        int8
	Direction   int8
	DisplayName pk.Option[chat.Message, *chat.Message]
}

type MapColorPatch struct {
	Columns uint8
	Rows    uint8
	X, Z    uint8
	Data    []pk.UnsignedByte
}

func (c *MapColorPatch) ReadFrom(r io.Reader) (n int64, err error) {
	c.Rows = 0
	c.X = 0
	c.Z = 0
	c.Data = nil
	t, err := (*pk.UnsignedByte)(&c.Columns).ReadFrom(r)
	n += t
	if err != nil {
		return n, err
	}
	if c.Columns <= 0 {
		return n, nil
	}
	for _, field := range []*uint8{&c.Rows, &c.X, &c.Z} {
		t, err = (*pk.UnsignedByte)(field).ReadFrom(r)
		n += t
		if err != nil {
			return n, err
		}
	}
	c.Data = nil
	t, err = pk.Array(&c.Data).ReadFrom(r)
	n += t
	return n, err
}

func (c MapColorPatch) WriteTo(w io.Writer) (n int64, err error) {
	t, err := pk.UnsignedByte(c.Columns).WriteTo(w)
	n += t
	if err != nil {
		return n, err
	}
	if c.Columns <= 0 {
		return n, nil
	}
	for _, value := range []uint8{c.Rows, c.X, c.Z} {
		t, err = pk.UnsignedByte(value).WriteTo(w)
		n += t
		if err != nil {
			return n, err
		}
	}
	t, err = pk.Array(&c.Data).WriteTo(w)
	n += t
	return n, err
}

//codec:gen
type MapData struct {
	MapID          int32 `mc:"VarInt"`
	Scale          int8
	Locked         bool
	HasDecorations bool
	Decorations    []MapIcon
	ColorPatch     MapColorPatch
}
