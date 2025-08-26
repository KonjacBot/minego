package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ConfigTransfer struct {
	Host string
	Port int32 `mc:"VarInt"`
}

func (*ConfigTransfer) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigTransfer
}

func init() {
	registerPacket(packetid.ClientboundConfigTransfer, func() ClientboundPacket {
		return &ConfigTransfer{}
	})
}
