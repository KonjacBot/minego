package server

import (
	"bytes"
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

	data := make([]byte, 32767)
	nn, err := io.ReadFull(r, data)
	if err != nil && (err != io.ErrUnexpectedEOF && err != io.EOF) {
		return n + int64(nn), err
	}
	p.Data = data
	return n + int64(nn), nil
}
