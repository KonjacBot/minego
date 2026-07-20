package server

import (
	"io"

	"github.com/KonjacBot/go-mc/data/packetid"
)

type ConfigAcceptCodeOfConduct struct{}

func (*ConfigAcceptCodeOfConduct) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundConfigAcceptCodeOfConduct
}

func (*ConfigAcceptCodeOfConduct) ReadFrom(io.Reader) (int64, error) {
	return 0, nil
}

func (ConfigAcceptCodeOfConduct) WriteTo(io.Writer) (int64, error) {
	return 0, nil
}

func init() {
	registerPacket(packetid.ServerboundConfigAcceptCodeOfConduct, func() ServerboundPacket {
		return &ConfigAcceptCodeOfConduct{}
	})
}
