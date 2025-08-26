package client

import (
	"io"

	"git.konjactw.dev/falloutBot/go-mc/chat"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
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
	t, err := (*pk.UnsignedByte)(&c.Columns).ReadFrom(r)
	if err != nil {
		return t, err
	}
	if c.Columns <= 0 {
		return t, err
	}
	a, err := (*pk.UnsignedByte)(&c.Rows).ReadFrom(r)
	b, err := (*pk.UnsignedByte)(&c.X).ReadFrom(r)
	d, err := (*pk.UnsignedByte)(&c.Z).ReadFrom(r)
	e, err := pk.Array(&c.Data).ReadFrom(r)
	return t + a + b + d + e, err
}

func (c MapColorPatch) WriteTo(w io.Writer) (n int64, err error) {
	n, err = pk.UnsignedByte(c.Columns).WriteTo(w)
	if c.Columns <= 0 {
		return n, err
	}
	n, err = pk.UnsignedByte(c.Rows).WriteTo(w)
	n, err = pk.UnsignedByte(c.X).WriteTo(w)
	n, err = pk.UnsignedByte(c.Z).WriteTo(w)
	n, err = pk.Array(&c.Data).WriteTo(w)
	return n, err
}

//codec:gen
type MapData struct {
	MapID  int32 `mc:"VarInt"`
	Scale  int8
	Locked bool
}
