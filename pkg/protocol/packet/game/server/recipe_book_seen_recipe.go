package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type RecipeBookSeenRecipe struct {
	RecipeID int32 `mc:"VarInt"`
}

func (*RecipeBookSeenRecipe) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundRecipeBookSeenRecipe
}

func init() {
	registerPacket(packetid.ServerboundRecipeBookSeenRecipe, func() ServerboundPacket {
		return &RecipeBookSeenRecipe{}
	})
}
