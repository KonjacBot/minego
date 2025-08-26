package server

import (
	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/yggdrasil/user"
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
