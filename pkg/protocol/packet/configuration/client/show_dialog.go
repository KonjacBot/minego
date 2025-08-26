package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
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
