package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type CommandSuggestion struct {
	TransactionID int32 `mc:"VarInt"`
	Text          string
}

func (CommandSuggestion) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundCommandSuggestion
}

func init() {
	registerPacket(packetid.ServerboundCommandSuggestion, func() ServerboundPacket {
		return &CommandSuggestion{}
	})
}
