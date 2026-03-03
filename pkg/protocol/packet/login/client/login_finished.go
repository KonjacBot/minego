package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/minego/pkg/protocol"
)

//codec:gen
type LoginLoginFinished struct {
	GameProfile protocol.GameProfile
}

func (*LoginLoginFinished) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginFinished
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginFinished, func() ClientboundPacket {
		return &LoginLoginFinished{}
	})
}
