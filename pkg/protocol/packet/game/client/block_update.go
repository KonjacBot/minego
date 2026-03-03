package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/net/packet"
)

// codec:gen
type BlockUpdate struct {
	Position   packet.Position
	BlockState int32 `mc:"VarInt"`
}

func (BlockUpdate) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockUpdate
}
