package client

import (
	"github.com/KonjacBot/go-mc/data/packetid"

	"github.com/KonjacBot/minego/v2/pkg/protocol"
)

//codec:gen
type KnownPack struct {
	Namespace string
	ID        string
	Version   string
}

//codec:gen
type ConfigSelectKnownPacks struct {
	KnownPacks []KnownPack
}

func (*ConfigSelectKnownPacks) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigSelectKnownPacks
}

func init() {
	registerPacket(packetid.ClientboundConfigSelectKnownPacks, func() ClientboundPacket {
		return &ConfigSelectKnownPacks{}
	})
}

// SupportedKnownPacks picks the packs from the server's offer that this
// protocol version ships with (only the vanilla core pack), which is what a
// client answers to ConfigSelectKnownPacks.
func SupportedKnownPacks(offered []KnownPack) []KnownPack {
	result := make([]KnownPack, 0, 1)
	for _, pack := range offered {
		if pack.Namespace == protocol.CorePackNamespace &&
			pack.ID == protocol.CorePackID &&
			pack.Version == protocol.VersionName {
			result = append(result, pack)
		}
	}
	return result
}
