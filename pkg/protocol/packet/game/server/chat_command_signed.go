package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type SignedSignatures struct {
	ArgumentName string
	Signature    []byte `mc:"ByteArray"`
}

//codec:gen
type ChatCommandSigned struct {
	Command            string
	Timestamp          int64
	Salt               int64
	ArgumentSignatures []SignedSignatures
	MessageCount       int32          `mc:"VarInt"`
	Acknowledged       pk.FixedBitSet `mc:"FixedBitSet" size:"20"`
	Checksum           int8
}

func (*ChatCommandSigned) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChatCommandSigned
}

func init() {
	registerPacket(packetid.ServerboundChatCommandSigned, func() ServerboundPacket {
		return &ChatCommandSigned{}
	})
}
