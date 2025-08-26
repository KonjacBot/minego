package client

import (
	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*AddEntity)(nil)
var _ packet.Field = (*AddEntity)(nil)

// AddEntityPacket
// codec:gen
type AddEntity struct {
	ID                              int32     `mc:"VarInt"`
	UUID                            uuid.UUID `mc:"UUID"`
	Type                            int32     `mc:"VarInt"`
	X, Y, Z                         float64
	XRot, YRot, YHeadRot            int8
	Data                            int32 `mc:"VarInt"`
	VelocityX, VelocityY, VelocityZ int16
}

func (AddEntity) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAddEntity
}
