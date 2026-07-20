package client

import (
	"io"

	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/nbt"
	"github.com/KonjacBot/go-mc/net/packet"
)

type ConfigShowDialog struct {
	DialogData nbt.RawMessage `mc:"NBT"`
}

func (*ConfigShowDialog) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigShowDialog
}

func (c *ConfigShowDialog) ReadFrom(r io.Reader) (n int64, err error) {
	return packet.NBT(&c.DialogData).ReadFrom(r)
}

func (c ConfigShowDialog) WriteTo(w io.Writer) (n int64, err error) {
	return packet.NBT(&c.DialogData).WriteTo(w)
}

func init() {
	registerPacket(packetid.ClientboundConfigShowDialog, func() ClientboundPacket {
		return &ConfigShowDialog{}
	})
}
