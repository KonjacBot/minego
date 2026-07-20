package server

import (
	"github.com/KonjacBot/go-mc/chat/sign"
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type Chat struct {
	Message      string
	Timestamp    int64
	Salt         int64
	HasSignature bool
	//opt:optional:HasSignature
	Signature    sign.Signature
	MessageCount int32          `mc:"VarInt"`
	Acknowledged pk.FixedBitSet `mc:"FixedBitSet" size:"20"`
	Checksum     uint8
}

func (*Chat) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChat
}

func init() {
	registerPacket(packetid.ServerboundChat, func() ServerboundPacket {
		return &Chat{}
	})
}
