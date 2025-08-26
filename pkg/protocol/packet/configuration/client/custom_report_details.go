package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"

	"git.konjactw.dev/patyhank/minego/pkg/protocol/packet/game/client"
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
