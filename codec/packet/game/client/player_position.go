package client

import (
	"github.com/Tnze/go-mc/data/packetid"
)

var _ ClientboundPacket = (*PlayerPosition)(nil)

//codec:gen
type PositionMoveRotation struct {
	X, Y, Z    float64
	YRot, XRot float32
}

//codec:gen
type PlayerPosition struct {
	ID        int32 `mc:"VarInt"`
	Change    PositionMoveRotation
	Relatives uint8
}

func (PlayerPosition) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerPosition
}
