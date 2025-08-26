package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type EditBook struct {
	Slot     int32 `mc:"VarInt"`
	Entries  []string
	HasTitle bool
	//opt:optional:HasTitle
	Title string
}

func (*EditBook) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundEditBook
}

func init() {
	registerPacket(packetid.ServerboundEditBook, func() ServerboundPacket {
		return &EditBook{}
	})
}
