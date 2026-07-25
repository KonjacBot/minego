package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/minego/pkg/protocol/wire"
)

//codec:gen
type ConfigDisconnect struct {
	Reason wire.Message
}

func (*ConfigDisconnect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigDisconnect
}

func init() {
	registerPacket(packetid.ClientboundConfigDisconnect, func() ClientboundPacket {
		return &ConfigDisconnect{}
	})
}
