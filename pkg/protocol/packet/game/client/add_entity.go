package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
	"github.com/google/uuid"
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
