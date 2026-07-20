package metadata

import (
	"io"

	"github.com/KonjacBot/go-mc/net/packet"
)

func (b *OptBlockState) ReadFrom(r io.Reader) (n int64, err error) {
	b.Value = 0
	var raw int32
	n, err = (*packet.VarInt)(&raw).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if raw != 0 {
		b.Value = raw
	}
	return n, nil
}

func (b OptBlockState) WriteTo(w io.Writer) (n int64, err error) {
	return (*packet.VarInt)(&b.Value).WriteTo(w)
}

func (b *OptVarInt) ReadFrom(r io.Reader) (n int64, err error) {
	b.Has = false
	b.Value = 0
	var raw int32
	n, err = (*packet.VarInt)(&raw).ReadFrom(r)
	if err != nil {
		return n, err
	}
	if raw != 0 {
		b.Has = true
		b.Value = raw - 1
	}
	return n, nil
}

func (b OptVarInt) WriteTo(w io.Writer) (n int64, err error) {
	if !b.Has {
		return packet.VarInt(0).WriteTo(w)
	}
	value := b.Value + 1
	return (*packet.VarInt)(&value).WriteTo(w)
}
