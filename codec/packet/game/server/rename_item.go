package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type RenameItem struct {
	ItemName string
}

func (RenameItem) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundRenameItem
}

func init() {
	registerPacket(packetid.ServerboundRenameItem, func() ServerboundPacket {
		return &RenameItem{}
	})
}
