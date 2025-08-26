package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
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
