package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/pkg/protocol/packet/game/client"
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
