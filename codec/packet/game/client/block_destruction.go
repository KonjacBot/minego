package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

var _ ClientboundPacket = (*BlockDestruction)(nil)

// codec:gen
type BlockDestruction struct {
	ID       int32 `mc:"VarInt"`
	Position packet.Position
	Progress uint8
}

func (BlockDestruction) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockDestruction
}
