package client

import (
	"io"

	pk "github.com/KonjacBot/go-mc/net/packet"
)

type StopSound struct {
	Flags  int8
	Source int32  `mc:"VarInt"`
	Sound  string `mc:"Identifier"`
}

func (s StopSound) WriteTo(w io.Writer) (n int64, err error) {
	temp, err := pk.Byte(s.Flags).WriteTo(w)
	n += temp
	if err != nil {
		return n, err
	}
	if s.Flags&0x01 != 0 {
		temp, err = pk.VarInt(s.Source).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	if s.Flags&0x02 != 0 {
		temp, err = pk.Identifier(s.Sound).WriteTo(w)
		n += temp
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (s *StopSound) ReadFrom(r io.Reader) (n int64, err error) {
	temp, err := (*pk.Byte)(&s.Flags).ReadFrom(r)
	n += temp
	if err != nil {
		return n, err
	}
	if s.Flags&0x01 != 0 {
		temp, err = (*pk.VarInt)(&s.Source).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	} else {
		s.Source = 0
	}
	if s.Flags&0x02 != 0 {
		temp, err = (*pk.Identifier)(&s.Sound).ReadFrom(r)
		n += temp
		if err != nil {
			return n, err
		}
	} else {
		s.Sound = ""
	}
	return n, nil
}
