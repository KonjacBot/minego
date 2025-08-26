package client

import (
	"io"

	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

type StopSound struct {
	Flags  int8
	Source int32  `mc:"VarInt"`
	Sound  string `mc:"Identifier"`
}

func (s StopSound) WriteTo(w io.Writer) (n int64, err error) {
	pk.Byte(s.Flags).WriteTo(w)
	if s.Flags&0x01 != 0 {
		pk.VarInt(s.Source).WriteTo(w)
	}
	if s.Flags&0x02 != 0 {
		pk.Identifier(s.Sound).WriteTo(w)
	}
	return
}

func (s *StopSound) ReadFrom(r io.Reader) (n int64, err error) {
	(*pk.Byte)(&s.Flags).ReadFrom(r)
	if s.Flags&0x01 != 0 {
		(*pk.VarInt)(&s.Source).ReadFrom(r)
	}
	if s.Flags&0x02 != 0 {
		(*pk.Identifier)(&s.Sound).ReadFrom(r)
	}
	return
}
