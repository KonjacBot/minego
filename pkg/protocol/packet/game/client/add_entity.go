package client

import (
	"github.com/KonjacBot/minego/pkg/protocol"
	"github.com/google/uuid"

	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/net/packet"
)

var _ ClientboundPacket = (*AddEntity)(nil)
var _ packet.Field = (*AddEntity)(nil)

// AddEntityPacket
// codec:gen
type AddEntity struct {
	ID                  int32     `mc:"VarInt"`
	UUID                uuid.UUID `mc:"UUID"`
	Type                int32     `mc:"VarInt"`
	X, Y, Z             float64
	Velocity            protocol.LpVec3
	Pitch, Yaw, HeadYaw packet.Angle
	Data                int32 `mc:"VarInt"`
}

func (AddEntity) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundAddEntity
}
