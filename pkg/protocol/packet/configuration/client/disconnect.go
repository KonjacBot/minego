package client

import (
	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

//codec:gen
type ConfigDisconnect struct {
	Reason chat.Message
}

func (*ConfigDisconnect) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigDisconnect
}

func init() {
	registerPacket(packetid.ClientboundConfigDisconnect, func() ClientboundPacket {
		return &ConfigDisconnect{}
	})
}
