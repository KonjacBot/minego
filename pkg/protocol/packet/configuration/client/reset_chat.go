package client

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type ConfigResetChat struct {
}

func (*ConfigResetChat) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigResetChat
}

func init() {
	registerPacket(packetid.ClientboundConfigResetChat, func() ClientboundPacket {
		return &ConfigResetChat{}
	})
}
