package server

import "github.com/Tnze/go-mc/data/packetid"

//codec:gen
type ChunkBatchReceived struct {
	ChunksPerTick float32
}

func (ChunkBatchReceived) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundChunkBatchReceived
}

func init() {
	registerPacket(packetid.ServerboundChunkBatchReceived, func() ServerboundPacket {
		return &ChunkBatchReceived{}
	})
}
