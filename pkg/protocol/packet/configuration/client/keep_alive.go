package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ConfigKeepAlive struct {
	ID int64
}

func (*ConfigKeepAlive) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigKeepAlive
}

func init() {
	registerPacket(packetid.ClientboundConfigKeepAlive, func() ClientboundPacket {
		return &ConfigKeepAlive{}
	})
}
