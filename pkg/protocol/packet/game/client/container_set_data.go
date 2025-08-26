package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/net/packet"
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
