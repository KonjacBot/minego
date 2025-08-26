package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
)

var _ ClientboundPacket = (*BlockChangedAck)(nil)

// codec:gen
type BlockChangedAck struct {
	Sequence int32 `mc:"VarInt"`
}

func (BlockChangedAck) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockChangedAck
}
