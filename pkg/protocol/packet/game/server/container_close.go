package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ContainerClose struct {
	WindowID int32 `mc:"VarInt"`
}

func (*ContainerClose) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundContainerClose
}

func init() {
	registerPacket(packetid.ServerboundContainerClose, func() ServerboundPacket {
		return &ContainerClose{}
	})
}
