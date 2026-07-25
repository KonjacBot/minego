package server

import (
	"fmt"

	"github.com/KonjacBot/go-mc/chat/sign"
	"github.com/KonjacBot/go-mc/data/packetid"
	pk "github.com/KonjacBot/go-mc/net/packet"
)

//codec:gen
type SignedSignatures struct {
	ArgumentName string
	Signature    sign.Signature
}

//codec:gen
type ChatCommandSigned struct {
	Command            string
	Timestamp          int64
	Salt               int64
	ArgumentSignatures []SignedSignatures
	MessageCount       int32          `mc:"VarInt"`
	Acknowledged       pk.FixedBitSet `mc:"FixedBitSet" size:"20"`
	Checksum           uint8
}

func (c ChatCommandSigned) Validate() error {
	if len(c.Acknowledged) != acknowledgedMessageBytes {
		return fmt.Errorf("acknowledged messages is %d bytes, want %d", len(c.Acknowledged), acknowledgedMessageBytes)
	}
	if c.MessageCount < 0 {
		return fmt.Errorf("message count less than zero")
	}
	return nil
}

func (*ChatCommandSigned) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChatCommandSigned
}

func init() {
	registerPacket(packetid.ServerboundChatCommandSigned, func() ServerboundPacket {
		return &ChatCommandSigned{}
	})
}
