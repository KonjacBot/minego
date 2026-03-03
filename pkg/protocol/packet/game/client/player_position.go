package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
)

var _ ClientboundPacket = (*PlayerPosition)(nil)

//codec:gen
type PlayerPosition struct {
	ID                              int32 `mc:"VarInt"`
	X, Y, Z                         float64
	VelocityX, VelocityY, VelocityZ float64
	XRot, YRot                      float32
	Flags                           int32
}

func (PlayerPosition) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundPlayerPosition
}
