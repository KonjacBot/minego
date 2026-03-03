package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
)

type ConfigShowDialog struct {
	client.ShowDialog
}

func (*ConfigShowDialog) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigShowDialog
}

func init() {
	registerPacket(packetid.ClientboundConfigShowDialog, func() ClientboundPacket {
		return &ConfigShowDialog{}
	})
}
