package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type SelectTrade struct {
	SelectedSlot int32 `mc:"VarInt"`
}

func (*SelectTrade) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundSelectTrade
}

func init() {
	registerPacket(packetid.ServerboundSelectTrade, func() ServerboundPacket {
		return &SelectTrade{}
	})
}
