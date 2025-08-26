package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type ChangeGameMode struct {
	GameMode int32 `mc:"VarInt"`
}

func (*ChangeGameMode) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChangeGameMode
}

func init() {
	registerPacket(packetid.ServerboundChangeGameMode, func() ServerboundPacket {
		return &ChangeGameMode{}
	})
}
