package client

import (
	"git.konjactw.dev/falloutBot/go-mc/data/packetid"
	"git.konjactw.dev/falloutBot/go-mc/nbt"
)

//codec:gen
type RegistryData struct {
	Entry   string `mc:"Identifier"`
	HasData bool
	//opt:optional:HasData
	Data nbt.RawMessage `mc:"NBT"`
}

//codec:gen
type ConfigRegistryData struct {
	RegistryID string `mc:"Identifier"`
	Data       []RegistryData
}

func (*ConfigRegistryData) PacketID() packetid.ClientboundPacketID {
	return packetid.ClientboundConfigRegistryData
}

func init() {
	registerPacket(packetid.ClientboundConfigRegistryData, func() ClientboundPacket {
		return &ConfigRegistryData{}
	})
}
