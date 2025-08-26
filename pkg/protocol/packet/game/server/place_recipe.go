package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type PlaceRecipe struct {
	WindowID int32 `mc:"VarInt"`
	RecipeID int32 `mc:"VarInt"`
	MakeAll  bool
}

func (*PlaceRecipe) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundPlaceRecipe
}

func init() {
	registerPacket(packetid.ServerboundPlaceRecipe, func() ServerboundPacket {
		return &PlaceRecipe{}
	})
}
