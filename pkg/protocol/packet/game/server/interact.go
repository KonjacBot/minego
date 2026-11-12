package server

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/v2/pkg/protocol"
)

//codec:gen
type Interact struct {
	EntityID             int32 `mc:"VarInt"`
	Hand                 int32 `mc:"VarInt"`
	Location             protocol.LpVec3
	UsingSecondaryAction bool
}

func (*Interact) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundInteract
}

func init() {
	registerPacket(packetid.ServerboundInteract, func() ServerboundPacket {
		return &Interact{}
	})
}
