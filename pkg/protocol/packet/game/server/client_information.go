package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ClientInformation struct {
	Location            string
	ViewDistance        int8
	ChatMode            int32 `mc:"VarInt"`
	ChatColor           bool
	DisplayedSkinParts  uint8
	MainHand            int32 `mc:"VarInt"`
	EnableTextFiltering bool
	AllowListing        bool
	ParticleStatus      int32 `mc:"VarInt"`
}

func (*ClientInformation) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundClientInformation
}

func init() {
	registerPacket(packetid.ServerboundClientInformation, func() ServerboundPacket {
		return &ClientInformation{}
	})
}
