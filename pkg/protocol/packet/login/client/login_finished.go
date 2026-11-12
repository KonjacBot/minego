package client

import (
	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/minego/v2/pkg/protocol"
)

//codec:gen
type LoginLoginFinished struct {
	GameProfile protocol.GameProfile
	SessionID   uuid.UUID `mc:"UUID"`
}

func (*LoginLoginFinished) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundLoginLoginFinished
}

func init() {
	registerPacket(packetid.ClientboundLoginLoginFinished, func() ClientboundPacket {
		return &LoginLoginFinished{}
	})
}
