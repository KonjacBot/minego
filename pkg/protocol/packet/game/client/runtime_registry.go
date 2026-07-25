package client

func init() {
	registerPacket(func() ClientboundPacket {
		return &AddResourcePack{}
	})
	registerPacket(func() ClientboundPacket {
		return &RemoveResourcePack{}
	})
}
