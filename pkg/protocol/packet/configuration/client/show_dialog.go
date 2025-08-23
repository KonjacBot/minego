package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
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
