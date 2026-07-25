package server

import (
	"fmt"

	"github.com/KonjacBot/go-mc/chat/sign"
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

const acknowledgedMessageBytes = 3

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

func (c Chat) Validate() error {
	if len(c.Acknowledged) != acknowledgedMessageBytes {
		return fmt.Errorf("acknowledged messages is %d bytes, want %d", len(c.Acknowledged), acknowledgedMessageBytes)
	}
	if c.MessageCount < 0 {
		return fmt.Errorf("message count less than zero")
	}
	return nil
}

func (*Chat) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChat
}

func init() {
	registerPacket(packetid.ServerboundChat, func() ServerboundPacket {
		return &Chat{}
	})
}
