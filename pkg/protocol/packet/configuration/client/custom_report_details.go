package client

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
	"github.com/Tnze/go-mc/data/packetid"
)

type ConfigCustomReportDetails struct {
	client.CustomReportDetails
}

func (*ConfigCustomReportDetails) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigCustomReportDetails
}

func init() {
	registerPacket(packetid.ClientboundConfigCustomReportDetails, func() ClientboundPacket {
		return &ConfigCustomReportDetails{}
	})
}
