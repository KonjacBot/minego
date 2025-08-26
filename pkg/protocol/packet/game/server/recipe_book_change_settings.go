package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type RecipeBookChangeSettings struct {
	BookId       int32 `mc:"VarInt"`
	BookOpen     bool
	FilterActive bool
}

func (*RecipeBookChangeSettings) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundRecipeBookChangeSettings
}

func init() {
	registerPacket(packetid.ServerboundRecipeBookChangeSettings, func() ServerboundPacket {
		return &RecipeBookChangeSettings{}
	})
}
