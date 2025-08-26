package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type AcceptTeleportation struct {
	TeleportID int32 `mc:"VarInt"`
}

func (*AcceptTeleportation) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundAcceptTeleportation
}

func init() {
	registerPacket(packetid.ServerboundAcceptTeleportation, func() ServerboundPacket {
		return &AcceptTeleportation{}
	})
}
