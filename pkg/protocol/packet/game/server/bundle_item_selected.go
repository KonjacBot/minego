package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type BundleItemSelected struct {
	SlotOfBundle int32 `mc:"VarInt"`
	SlotInBundle int32 `mc:"VarInt"`
}

func (*BundleItemSelected) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundBundleItemSelected
}

func init() {
	registerPacket(packetid.ServerboundBundleItemSelected, func() ServerboundPacket {
		return &BundleItemSelected{}
	})
}
