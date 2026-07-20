package client

import (
	"errors"
	"io"

	"github.com/KonjacBot/go-mc/chat/sign"
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*DeleteChat)(nil)
var _ packet.Field = (*DeleteChat)(nil)

// DeleteChatPacket
//
//codec:gen
type DeleteChat struct {
	MessageSignature PackedMessageSignature
}

// PackedMessageSignature contains either a cached signature ID or an inline
// 256-byte signature. Cached IDs are zero-based; zero on the wire means inline.
type PackedMessageSignature struct {
	ID        int32
	Signature *sign.Signature
}

func (p *PackedMessageSignature) ReadFrom(r io.Reader) (n int64, err error) {
	var encodedID packet.VarInt
	n, err = encodedID.ReadFrom(r)
	if err != nil {
		return n, err
	}

	p.ID = int32(encodedID) - 1
	if p.ID >= 0 {
		p.Signature = nil
		return n, nil
	}
	if p.ID != -1 {
		return n, errors.New("packed signature id less than zero")
	}

	p.Signature = new(sign.Signature)
	nn, err := io.ReadFull(r, p.Signature[:])
	return n + int64(nn), err
}

func (p PackedMessageSignature) WriteTo(w io.Writer) (n int64, err error) {
	if p.Signature != nil {
		n, err = packet.VarInt(0).WriteTo(w)
		if err != nil {
			return n, err
		}
		nn, err := w.Write(p.Signature[:])
		return n + int64(nn), err
	}
	if p.ID < 0 {
		return 0, errors.New("cached signature id less than zero")
	}
	return packet.VarInt(p.ID + 1).WriteTo(w)
}

func (DeleteChat) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundDeleteChat
}
