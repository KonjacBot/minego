package client

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

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
