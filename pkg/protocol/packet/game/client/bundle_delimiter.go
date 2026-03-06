package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"
)

// BundleDelimiter
//
// codec:gen
type BundleDelimiter struct {
}

func (BundleDelimiter) PacketID() packetid.ClientboundPacketID {
	return packetid.BundleDelimiter
}
