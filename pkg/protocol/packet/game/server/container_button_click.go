package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ContainerButtonClick struct {
	WindowID int32 `mc:"VarInt"`
	ButtonID int32 `mc:"VarInt"`
}

func (*ContainerButtonClick) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundContainerButtonClick
}

func init() {
	registerPacket(packetid.ServerboundContainerButtonClick, func() ServerboundPacket {
		return &ContainerButtonClick{}
	})
}
