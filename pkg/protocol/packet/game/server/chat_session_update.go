package server

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/yggdrasil/user"
	"github.com/google/uuid"
)

//codec:gen
type ChatSessionUpdate struct {
	SessionId uuid.UUID `mc:"UUID"`
	PublicKey user.PublicKey
}

func (*ChatSessionUpdate) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChatSessionUpdate
}

func init() {
	registerPacket(packetid.ServerboundChatSessionUpdate, func() ServerboundPacket {
		return &ChatSessionUpdate{}
	})
}
