package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

// codec:gen
type BlockUpdate struct {
	Position   packet.Position
	BlockState int32 `mc:"VarInt"`
}

func (BlockUpdate) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundBlockUpdate
}
