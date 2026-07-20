package server

import (
	"bytes"
	"errors"
	"io"

	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

type CustomPayload struct {
	Channel string
	Data    []byte
}

func (*CustomPayload) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundCustomPayload
}

func init() {
	registerPacket(packetid.ServerboundCustomPayload, func() ServerboundPacket {
		return &CustomPayload{}
	})
}

func (p CustomPayload) WriteTo(w io.Writer) (n int64, err error) {
	if len(p.Data) > 32767 {
		return 0, errors.New("custom payload exceeds 32767 bytes")
	}
	n, err = pk.Identifier(p.Channel).WriteTo(w)
	if err != nil {
		return n, err
	}

	nn, err := bytes.NewBuffer(p.Data).WriteTo(w)
	return n + nn, err
}

func (p *CustomPayload) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = (*pk.Identifier)(&p.Channel).ReadFrom(r)
	if err != nil {
		return
	}

	data, err := io.ReadAll(io.LimitReader(r, 32768))
	if err != nil {
		return n, err
	}
	if len(data) > 32767 {
		return n + int64(len(data)), errors.New("custom payload exceeds 32767 bytes")
	}
	p.Data = data
	return n + int64(len(data)), nil
}
