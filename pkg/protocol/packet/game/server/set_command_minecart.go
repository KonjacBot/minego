package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type SetCommandMinecart struct {
	EntityID    int32 `mc:"VarInt"`
	Command     string
	TrackOutput bool
}

func (*SetCommandMinecart) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSetCommandMinecart
}

func init() {
	registerPacket(packetid.ServerboundSetCommandMinecart, func() ServerboundPacket {
		return &SetCommandMinecart{}
	})
}
