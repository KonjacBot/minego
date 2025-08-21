package client

import (
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net/packet"
)

var _ ClientboundPacket = (*ContainerSetData)(nil)
var _ packet.Field = (*ContainerSetData)(nil)

// ContainerSetDataPacket
//
//codec:gen
type ContainerSetData struct {
	ContainerID int8
	ID          int16
	Value       int16
}

func (ContainerSetData) ClientboundPacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundContainerSetData
}
