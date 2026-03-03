package client

import "github.com/KonjacBot/go-mc/data/packetid"

//codec:gen
type ConfigClearDialog struct {
}

func (*ConfigClearDialog) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigClearDialog
}

func init() {
	registerPacket(packetid.ClientboundConfigClearDialog, func() ClientboundPacket {
		return &ConfigClearDialog{}
	})
}
