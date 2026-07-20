package client

import "github.com/KonjacBot/go-mc/data/packetid"

func init() {
	ClientboundPackets[packetid.ClientboundResourcePackPush] = func() ClientboundPacket {
		return &AddResourcePack{}
	}
	ClientboundPackets[packetid.ClientboundResourcePackPop] = func() ClientboundPacket {
		return &RemoveResourcePack{}
	}
}
