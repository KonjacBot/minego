package server

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
)

//codec:gen
type Chat struct {
	Message      string
	Timestamp    int64
	Salt         int64
	HasSignature bool
	//opt:optional:HasSignature
	Signature    []byte         `mc:"ByteArray"`
	MessageCount int32          `mc:"VarInt"`
	Acknowledged pk.FixedBitSet `mc:"FixedBitSet" size:"20"`
	Checksum     int8
}

func (*Chat) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChat
}

func init() {
	registerPacket(packetid.ServerboundChat, func() ServerboundPacket {
		return &Chat{}
	})
}
