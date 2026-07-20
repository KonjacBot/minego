package client

import (
	"io"

	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/codecutil"
)

type ConfigCodeOfConduct struct {
	CodeOfConduct string
}

func (*ConfigCodeOfConduct) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigCodeOfConduct
}

func (c *ConfigCodeOfConduct) ReadFrom(r io.Reader) (n int64, err error) {
	return codecutil.BoundedString{Value: &c.CodeOfConduct, MaxChars: 32767}.ReadFrom(r)
}

func (c ConfigCodeOfConduct) WriteTo(w io.Writer) (n int64, err error) {
	return codecutil.BoundedString{Value: &c.CodeOfConduct, MaxChars: 32767}.WriteTo(w)
}

func init() {
	registerPacket(packetid.ClientboundConfigCodeOfConduct, func() ClientboundPacket {
		return &ConfigCodeOfConduct{}
	})
}
