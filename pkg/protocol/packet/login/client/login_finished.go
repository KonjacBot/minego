package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/patyhank/minego/pkg/protocol"
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
